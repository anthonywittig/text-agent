package main

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/ttacon/libphonenumber"
)

// https://docs.aws.amazon.com/bedrock/latest/userguide/agents-lambda.html
type AgentRequest struct {
	MessageVersion string `json:"messageVersion"`
	Function       string `json:"function"`
	Parameters     []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"parameters"`
	InputText string `json:"inputText"`
	SessionId string `json:"sessionId"`
	Agent     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Id      string `json:"id"`
		Alias   string `json:"alias"`
	} `json:"agent"`
	ActionGroup string `json:"actionGroup"`
	// Not sure what this looks like in practice.
	SessionAttributes interface{} `json:"sessionAttributes,omitempty"`
	// Not sure what this looks like in practice.
	PromptSessionAttributes interface{} `json:"promptSessionAttributes,omitempty"`
}

// https://docs.aws.amazon.com/bedrock/latest/userguide/agents-lambda.html
type AgentResponse struct {
	MessageVersion string                `json:"messageVersion"`
	Response       AgentResponseResponse `json:"response"`
	// SessionAttributes           interface{} `json:"sessionAttributes,omitempty"`
	// PromptSessionAttributes     interface{} `json:"promptSessionAttributes,omitempty"`
	// KnowledgeBasesConfiguration []struct {
	// 	KnowledgeBaseId        string `json:"knowledgeBaseId"`
	// 	RetrievalConfiguration struct {
	// 		VectorSearchConfiguration struct {
	// 			NumberOfResults int `json:"numberOfResults"`
	// 			Filter          struct {
	// 				RetrievalFilter struct {
	// 					Field    string `json:"field"`
	// 					Operator string `json:"operator"`
	// 					Value    string `json:"value"`
	// 				} `json:"filter"`
	// 			} `json:"vectorSearchConfiguration"`
	// 		} `json:"retrievalConfiguration"`
	// 	} `json:"retrievalConfiguration"`
	// } `json:"knowledgeBasesConfiguration"`
}

type AgentResponseResponse struct {
	ActionGroup      string                                `json:"actionGroup"`
	Function         string                                `json:"function"`
	FunctionResponse AgentResponseResponseFunctionResponse `json:"functionResponse"`
}

type AgentResponseResponseFunctionResponse struct {
	ResponseState string                                            `json:"responseState"`
	ResponseBody  AgentResponseResponseFunctionResponseResponseBody `json:"responseBody"`
}

type AgentResponseResponseFunctionResponseResponseBody struct {
	ContentType AgentResponseResponseFunctionResponseResponseBodyContentType `json:"TEXT"`
}

type AgentResponseResponseFunctionResponseResponseBodyContentType struct {
	Body string `json:"body"` // This should be a JSON string.
}

const (
	tableName = "text-agent-task-tracking"
)

func handleRequest(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	repo, err := task_repository.New(ctx, tableName)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create repository")
		return AgentResponse{
			MessageVersion: "1.0",
			Response: AgentResponseResponse{
				ActionGroup: payload.ActionGroup,
				Function:    payload.Function,
				FunctionResponse: AgentResponseResponseFunctionResponse{
					ResponseState: "FAILURE",
					ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
						ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
							Body: "{\"message\": \"Internal error\"}",
						},
					},
				},
			},
		}, nil
	}

	// Iterate over the parameters to find where `name` is "phone_numbers"
	phoneNumbers := []string{}
	for _, param := range payload.Parameters {
		if param.Name == "phone_numbers" {
			// The value will be a string like "[5551112222, 5551113333]".
			phoneNumbers = strings.Split(strings.Trim(param.Value, "[]"), ",")
			break
		}
	}
	if len(phoneNumbers) == 0 {
		logger.Error().Msg("no phone numbers found")
		return AgentResponse{
			MessageVersion: "1.0",
			Response: AgentResponseResponse{
				ActionGroup: payload.ActionGroup,
				Function:    payload.Function,
				FunctionResponse: AgentResponseResponseFunctionResponse{
					ResponseState: "FAILURE",
					ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
						ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
							Body: "{\"message\": \"No phone numbers found\"}",
						},
					},
				},
			},
		}, nil
	}

	e164PhoneNumbers := make([]string, len(phoneNumbers))
	for i, phoneNumber := range phoneNumbers {
		number, err := libphonenumber.Parse(phoneNumber, "US")
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse phone number")
			return AgentResponse{
				MessageVersion: "1.0",
				Response: AgentResponseResponse{
					ActionGroup: payload.ActionGroup,
					Function:    payload.Function,
					FunctionResponse: AgentResponseResponseFunctionResponse{
						ResponseState: "FAILURE",
						ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
							ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
								Body: "{\"message\": \"Unable to parse phone number " + phoneNumber + "\"}",
							},
						},
					},
				},
			}, nil
		}
		e164PhoneNumbers[i] = libphonenumber.Format(number, libphonenumber.E164)
	}

	sort.Strings(e164PhoneNumbers)
	conversationID := strings.Join(e164PhoneNumbers, "_")

	logger.Info().Str("conversation_id", conversationID).Msg("Processing conversation")

	tasks, err := repo.ListTasksByConversation(conversationID)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to list tasks")
		return AgentResponse{
			MessageVersion: "1.0",
			Response: AgentResponseResponse{
				ActionGroup: payload.ActionGroup,
				Function:    payload.Function,
				FunctionResponse: AgentResponseResponseFunctionResponse{
					ResponseState: "FAILURE",
					ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
						ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
							Body: "{\"message\": \"Internal error\"}",
						},
					},
				},
			},
		}, nil
	}

	taskString, err := json.Marshal(tasks)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal tasks")
	}

	return AgentResponse{
		MessageVersion: "1.0",
		Response: AgentResponseResponse{
			ActionGroup: payload.ActionGroup,
			Function:    payload.Function,
			FunctionResponse: AgentResponseResponseFunctionResponse{
				ResponseState: "SUCCESS",
				ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
					ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
						Body: string(taskString),
					},
				},
			},
		},
	}, nil
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	requestWrapper := func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		requestID := "unknown"
		if lc != nil {
			requestID = lc.AwsRequestID
		}
		logger := logger.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Interface("request", payload).Msg("received request")

		var request AgentRequest
		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal request")
			response := AgentResponse{
				MessageVersion: "1.0",
				Response: AgentResponseResponse{
					ActionGroup: "invalid_request",
					Function:    "invalid_request",
					FunctionResponse: AgentResponseResponseFunctionResponse{
						ResponseState: "FAILURE",
						ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
							ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
								Body: "{\"message\": \"Invalid request\"}",
							},
						},
					},
				},
			}
			responseJSON, _ := json.Marshal(response)
			return responseJSON, nil
		}

		logger.Info().Interface("request", request).Msg("parsed request")

		response, err := handleRequest(ctx, request)
		if err != nil {
			logger.Error().Err(err).Msg("failed to handle request")
			return nil, err
		}

		responseJson, _ := json.Marshal(response)
		logger.Info().Interface("response", response).Msg("sending response")
		return responseJson, nil
	}

	lambda.Start(requestWrapper)
}

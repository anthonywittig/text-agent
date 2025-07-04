package agent_action_consumer

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/rs/zerolog"
	"github.com/ttacon/libphonenumber"
)

type Consumer struct {
	repo task_repository.TaskRepository
}

func NewConsumer(repo task_repository.TaskRepository) *Consumer {
	return &Consumer{repo: repo}
}

func (c *Consumer) HandleRequest(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	switch payload.ActionGroup {
	case "TaskTrackingList":
		return c.handleTaskTrackingList(ctx, payload)
	default:
		logger.Error().Str("action_group", payload.ActionGroup).Msg("unknown action group")
		return AgentResponse{
			MessageVersion: "1.0",
			Response: AgentResponseResponse{
				ActionGroup: payload.ActionGroup,
				Function:    payload.Function,
				FunctionResponse: AgentResponseResponseFunctionResponse{
					ResponseState: "FAILURE",
					ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
						ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
							Body: "{\"message\": \"Unknown action group\"}",
						},
					},
				},
			},
		}, nil
	}
}

func (c *Consumer) handleTaskTrackingList(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

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

	tasks, err := c.repo.ListTasksByConversation(conversationID)
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
				ResponseState: "REPROMPT",
				ResponseBody: AgentResponseResponseFunctionResponseResponseBody{
					ContentType: AgentResponseResponseFunctionResponseResponseBodyContentType{
						Body: string(taskString),
					},
				},
			},
		},
	}, nil
}

package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/agent_action_consumer"
	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
)

const (
	tableName = "text-agent-task-tracking"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctx := context.Background()

	repo, err := task_repository.New(ctx, tableName)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create repository")
	}

	consumer := agent_action_consumer.NewConsumer(repo)

	requestWrapper := func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		requestID := "unknown"
		if lc != nil {
			requestID = lc.AwsRequestID
		}
		logger := logger.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Interface("request", payload).Msg("received request")

		var request agent_action_consumer.AgentRequest
		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal request")
			response := agent_action_consumer.AgentResponse{
				MessageVersion: "1.0",
				Response: agent_action_consumer.AgentResponseResponse{
					ActionGroup: "invalid_request",
					Function:    "invalid_request",
					FunctionResponse: agent_action_consumer.AgentResponseResponseFunctionResponse{
						ResponseState: "FAILURE",
						ResponseBody: agent_action_consumer.AgentResponseResponseFunctionResponseResponseBody{
							ContentType: agent_action_consumer.AgentResponseResponseFunctionResponseResponseBodyContentType{
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

		response, err := consumer.HandleRequest(ctx, request)
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

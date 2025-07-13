package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/agent_action_consumer"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/agent_service"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/message_repository"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/secrets_service"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctx := context.Background()

	agentAliasIdSecretId := os.Getenv("AGENT_ALIAS_ID_SECRET_ID")
	if agentAliasIdSecretId == "" {
		logger.Fatal().Msg("AGENT_ALIAS_ID_SECRET_ID is not set")
	}

	agentIdSecretId := os.Getenv("AGENT_ID_SECRET_ID")
	if agentIdSecretId == "" {
		logger.Fatal().Msg("AGENT_ID_SECRET_ID is not set")
	}

	secretsService, err := secrets_service.NewAwsSecretsService(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create secrets service")
	}

	agentAliasId, err := secretsService.GetSecret(ctx, agentAliasIdSecretId)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get agent alias ID")
	}

	agentId, err := secretsService.GetSecret(ctx, agentIdSecretId)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get agent ID")
	}

	agentService, err := agent_service.NewAws(ctx, agentAliasId, agentId)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create agent service")
	}

	repo, err := message_repository.New(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create repository")
	}

	consumer := agent_action_consumer.NewConsumer(agentService, repo)

	requestWrapper := func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		requestID := "unknown"
		if lc != nil {
			requestID = lc.AwsRequestID
		}
		logger := logger.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Interface("request", payload).Msg("received request")

		var request types.AgentRequest
		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal request")
			response := types.AgentResponse{
				MessageVersion: "1.0",
				Response: types.AgentResponseResponse{
					ActionGroup: "invalid_request",
					Function:    "invalid_request",
					FunctionResponse: types.AgentResponseResponseFunctionResponse{
						ResponseState: "FAILURE",
						ResponseBody: types.AgentResponseResponseFunctionResponseResponseBody{
							ContentType: types.AgentResponseResponseFunctionResponseResponseBodyContentType{
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

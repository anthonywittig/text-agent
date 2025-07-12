package agent_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Aws struct {
	agentAliasId string
	agentId      string
	bedrockAgent *bedrockagentruntime.Client
}

func NewAws(ctx context.Context, agentAliasId string, agentId string) (AgentService, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %v", err)
	}

	bedrockAgent := bedrockagentruntime.NewFromConfig(cfg)
	return &Aws{
		agentAliasId: agentAliasId,
		agentId:      agentId,
		bedrockAgent: bedrockAgent,
	}, nil
}

func (a *Aws) InvokeAgent(ctx context.Context, input string) (string, error) {
	logger := zerolog.Ctx(ctx)

	streamingConfigurations := types.StreamingConfigurations{
		StreamFinalResponse: true,
	}
	sessionId := uuid.New().String()
	invokeInput := &bedrockagentruntime.InvokeAgentInput{
		AgentAliasId:            &a.agentAliasId,
		AgentId:                 &a.agentId,
		InputText:               &input,
		SessionId:               &sessionId,
		EnableTrace:             aws.Bool(true),
		StreamingConfigurations: &streamingConfigurations,
	}

	logger.Info().
		Str("agentAliasId", a.agentAliasId).
		Str("agentId", a.agentId).
		Str("sessionId", sessionId).
		Str("input", input).
		Interface("streamingConfig", streamingConfigurations).
		Msg("invoking agent")

	invokeOutput, err := a.bedrockAgent.InvokeAgent(ctx, invokeInput)
	if err != nil {
		return "", fmt.Errorf("failed to invoke agent: %w", err)
	}

	logger.Info().
		Interface("resultMetadata", invokeOutput.ResultMetadata).
		Interface("contentType", invokeOutput.ContentType).
		Msg("agent invocation successful")

	stream := invokeOutput.GetStream()

	var response string
	logger.Info().Msg("starting to process stream events")

	for event := range stream.Events() {
		logger.Debug().Interface("event", event).Msg("received stream event")

		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			if e.Value.Bytes != nil {
				chunkText := string(e.Value.Bytes)
				response += chunkText
				logger.Debug().Str("chunk", chunkText).Msg("added chunk to response")
			}
		case *types.ResponseStreamMemberTrace:
			logger.Debug().Interface("trace", e.Value).Msg("received trace event")
		case *types.ResponseStreamMemberReturnControl:
			logger.Debug().Interface("returnControl", e.Value).Msg("received return control event")
		default:
			logger.Debug().Interface("event", event).Msg("received unknown event type")
		}
	}

	logger.Info().Str("finalResponse", response).Int("responseLength", len(response)).Msg("finished processing stream")

	// If we got no response, log a warning
	if response == "" {
		logger.Warn().Msg("agent returned empty response - this might indicate the agent didn't understand the input or couldn't take action")
	}

	return response, nil
}

package agent_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"

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

func (a *Aws) InvokeAgent(ctx context.Context, input string) error {
	logger := zerolog.Ctx(ctx)

	streamingConfigurations := awsTypes.StreamingConfigurations{
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
		return fmt.Errorf("failed to invoke agent: %w", err)
	}

	logger.Info().
		Interface("resultMetadata", invokeOutput.ResultMetadata).
		Interface("contentType", invokeOutput.ContentType).
		Msg("agent invocation successful")

	stream := invokeOutput.GetStream()

	logger.Info().Msg("starting to process stream events")

	for event := range stream.Events() {
		switch e := event.(type) {
		case *awsTypes.ResponseStreamMemberChunk:
			logger.Debug().Str("chunk", string(e.Value.Bytes)).Msg("received chunk")
		case *awsTypes.ResponseStreamMemberTrace:
			traceBytes, err := json.Marshal(e.Value)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal trace")
			}
			var trace types.AgentTrace
			err = json.Unmarshal(traceBytes, &trace)
			if err != nil {
				logger.Error().Err(err).Msg("failed to unmarshal trace")
			}

			var traceText types.AgentTraceTextFromJson
			err = json.Unmarshal([]byte(trace.Trace.Value.Value.Text), &traceText)
			if err != nil {
				logger.Error().Err(err).Msg("failed to unmarshal trace text")
			}

			for _, message := range traceText.Messages {
				logger.Debug().Str("message", message.Content).Str("role", message.Role).Msg("trace message")
			}

		case *awsTypes.ResponseStreamMemberReturnControl:
			logger.Debug().Interface("returnControl", e.Value).Msg("received return control event")
		default:
			logger.Warn().Interface("event", event).Msg("received unknown event type")
		}
	}

	return nil
}

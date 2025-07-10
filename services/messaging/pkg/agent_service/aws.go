package agent_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/google/uuid"
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
	sessionId := uuid.New().String()
	invokeInput := &bedrockagentruntime.InvokeAgentInput{
		AgentAliasId: &a.agentAliasId,
		AgentId:      &a.agentId,
		InputText:    &input,
		SessionId:    &sessionId,
	}

	invokeOutput, err := a.bedrockAgent.InvokeAgent(ctx, invokeInput)
	if err != nil {
		return "", err
	}

	stream := invokeOutput.GetStream()

	var response string
	for event := range stream.Events() {
		if chunk, ok := event.(*types.ResponseStreamMemberChunk); ok && chunk.Value.Bytes != nil {
			response += string(chunk.Value.Bytes)
		}
	}

	return response, nil
}

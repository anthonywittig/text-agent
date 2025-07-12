package agent_action_consumer

import (
	"context"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/agent_service"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/message_repository"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/types"
	"github.com/rs/zerolog"
)

type Consumer struct {
	agentService agent_service.AgentService
	repo         message_repository.MessageRepository
}

func NewConsumer(agentService agent_service.AgentService, repo message_repository.MessageRepository) *Consumer {
	return &Consumer{agentService: agentService, repo: repo}
}

func (c *Consumer) HandleRequest(ctx context.Context, payload types.AgentRequest) (types.AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	switch payload.Function {
	case "messaging_create":
		return c.handleMessageCreate(ctx, payload)
	case "messaging_list_recent":
		return c.handleMessageListRecent(ctx, payload)
	default:
		logger.Error().Str("function", payload.Function).Msg("unknown function")
		return types.AgentResponse{
			MessageVersion: "1.0",
			Response: types.AgentResponseResponse{
				ActionGroup: payload.ActionGroup,
				Function:    payload.Function,
				FunctionResponse: types.AgentResponseResponseFunctionResponse{
					ResponseState: "FAILURE",
					ResponseBody: types.AgentResponseResponseFunctionResponseResponseBody{
						ContentType: types.AgentResponseResponseFunctionResponseResponseBodyContentType{
							Body: "{\"message\": \"Unknown action group\"}",
						},
					},
				},
			},
		}, nil
	}
}

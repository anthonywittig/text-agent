package agent_action_consumer

import (
	"context"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/message_repository"
	"github.com/rs/zerolog"
)

type Consumer struct {
	repo message_repository.MessageRepository
}

func NewConsumer(repo message_repository.MessageRepository) *Consumer {
	return &Consumer{repo: repo}
}

func (c *Consumer) HandleRequest(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	switch payload.Function {
	case "message_create":
		return c.handleMessageCreate(ctx, payload)
	case "message_list_recent":
		return c.handleMessageListRecent(ctx, payload)
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

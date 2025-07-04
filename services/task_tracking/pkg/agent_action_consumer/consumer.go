package agent_action_consumer

import (
	"context"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/rs/zerolog"
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
	case "TaskTrackingCreate":
		return c.handleTaskTrackingCreate(ctx, payload)
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

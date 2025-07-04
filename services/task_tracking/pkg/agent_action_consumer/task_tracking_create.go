package agent_action_consumer

import (
	"context"
	"encoding/json"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/rs/zerolog"
)

type TaskTrackingCreateResponse struct {
	Message string                `json:"message"`
	Task    *task_repository.Task `json:"task"`
}

func (c *Consumer) handleTaskTrackingCreate(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Interface("payload", payload).Msg("handleTaskTrackingCreate")

	conversationId, err := getConversationId(ctx, payload)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	task, err := c.repo.CreateTask(
		conversationId,
		getParameter(payload, "name"),
		getParameter(payload, "description"),
		getParameter(payload, "source"),
	)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	response := TaskTrackingCreateResponse{
		Message: "Task created successfully",
		Task:    task,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
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
						Body: string(responseJson),
					},
				},
			},
		},
	}, nil
}

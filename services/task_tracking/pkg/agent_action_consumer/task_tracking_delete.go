package agent_action_consumer

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

type TaskTrackingDeleteResponse struct {
	Message string `json:"message"`
}

func (c *Consumer) handleTaskTrackingDelete(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Interface("payload", payload).Msg("handleTaskTrackingDelete")

	err := c.repo.DeleteTask(
		getParameter(payload, "task_id"),
	)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	response := TaskTrackingDeleteResponse{
		Message: "Task deleted successfully",
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

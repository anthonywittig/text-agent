package agent_action_consumer

import (
	"context"
	"encoding/json"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/message_repository"
	"github.com/rs/zerolog"
)

type MessageCreateResponse struct {
	Info    string                      `json:"info"`
	Message *message_repository.Message `json:"message"`
}

func (c *Consumer) handleMessageCreate(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Interface("payload", payload).Msg("handleMessageCreate")

	conversationId, err := getConversationId(ctx, payload)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	message, err := c.repo.CreateMessage(
		conversationId,
		getParameter(payload, "from"),
		getParameter(payload, "body"),
	)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	response := MessageCreateResponse{
		Info:    "Message created successfully",
		Message: message,
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

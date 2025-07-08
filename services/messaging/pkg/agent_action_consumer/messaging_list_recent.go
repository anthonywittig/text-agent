package agent_action_consumer

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

func (c *Consumer) handleMessageListRecent(ctx context.Context, payload AgentRequest) (AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	conversationId, err := getConversationId(ctx, payload)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	logger.Info().Str("conversation_id", conversationId).Msg("Processing conversation")

	messages, err := c.repo.ListRecentMessagesByConversation(conversationId)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationId).Msg("Failed to list messages")
		return getFailureResponse(payload, "Internal error"), nil
	}

	messageString, err := json.Marshal(messages)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal messages")
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
						Body: string(messageString),
					},
				},
			},
		},
	}, nil
}

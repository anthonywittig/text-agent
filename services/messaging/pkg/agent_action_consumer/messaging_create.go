package agent_action_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/message_repository"
	"github.com/anthonywittig/text-agent/services/messaging/pkg/types"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog"
)

type MessageCreateResponse struct {
	Info    string                      `json:"info"`
	Message *message_repository.Message `json:"message"`
}

func (c *Consumer) handleMessageCreate(ctx context.Context, payload types.AgentRequest) (types.AgentResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Interface("payload", payload).Msg("handleMessageCreate")

	conversationId, err := getConversationId(ctx, payload)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	// From can be things like "Assistant" or a phone number.
	from := getParameter(payload, "from")
	parsedFrom, err := phonenumbers.Parse(from, "US") // We assume US for now.
	if err == nil {
		from = phonenumbers.Format(parsedFrom, phonenumbers.E164)
	}

	message, err := c.repo.CreateMessage(
		conversationId,
		from,
		getParameter(payload, "body"),
	)
	if err != nil {
		return getFailureResponse(payload, err.Error()), nil
	}

	err = c.invokeAgent(ctx, payload)
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

	return types.AgentResponse{
		MessageVersion: "1.0",
		Response: types.AgentResponseResponse{
			ActionGroup: payload.ActionGroup,
			Function:    payload.Function,
			FunctionResponse: types.AgentResponseResponseFunctionResponse{
				ResponseState: "REPROMPT",
				ResponseBody: types.AgentResponseResponseFunctionResponseResponseBody{
					ContentType: types.AgentResponseResponseFunctionResponseResponseBodyContentType{
						Body: string(responseJson),
					},
				},
			},
		},
	}, nil
}

func (c *Consumer) invokeAgent(ctx context.Context, payload types.AgentRequest) error {
	logger := zerolog.Ctx(ctx)

	// If the message is from an agent, don't do anything.
	if payload.Agent.Name != "" {
		logger.Info().Str("agent_name", payload.Agent.Name).Msg("message from agent, skipping")
		return nil
	}

	err := c.agentService.InvokeAgent(ctx, "A new message was received for the conversation between these numbers: "+getParameter(payload, "conversation_phone_numbers"))
	if err != nil {
		return fmt.Errorf("failed to invoke agent: %w", err)
	}

	return nil
}

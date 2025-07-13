package agent_action_consumer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/anthonywittig/text-agent/services/messaging/pkg/types"
	"github.com/rs/zerolog"
	"github.com/ttacon/libphonenumber"
)

func getConversationId(ctx context.Context, payload types.AgentRequest) (string, error) {
	logger := zerolog.Ctx(ctx)

	phoneNumbers := strings.Split(strings.Trim(getParameter(payload, "conversation_phone_numbers"), "[]"), ",")
	if len(phoneNumbers) == 0 {
		logger.Error().Msg("no phone numbers found")
		return "", errors.New("no phone numbers found")
	}

	e164PhoneNumbers := make([]string, len(phoneNumbers))
	for i, phoneNumber := range phoneNumbers {
		number, err := libphonenumber.Parse(phoneNumber, "US")
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse phone number")
			return "", errors.New("failed to parse phone number " + phoneNumber)
		}
		e164PhoneNumbers[i] = libphonenumber.Format(number, libphonenumber.E164)
	}

	sort.Strings(e164PhoneNumbers)
	conversationId := strings.Join(e164PhoneNumbers, "_")
	return conversationId, nil
}

func getFailureResponse(payload types.AgentRequest, message string) types.AgentResponse {
	return types.AgentResponse{
		MessageVersion: "1.0",
		Response: types.AgentResponseResponse{
			ActionGroup: payload.ActionGroup,
			Function:    payload.Function,
			FunctionResponse: types.AgentResponseResponseFunctionResponse{
				ResponseState: "FAILURE",
				ResponseBody: types.AgentResponseResponseFunctionResponseResponseBody{
					ContentType: types.AgentResponseResponseFunctionResponseResponseBodyContentType{
						Body: fmt.Sprintf("{\"message\": \"%s\"}", message),
					},
				},
			},
		},
	}
}

func getParameter(payload types.AgentRequest, name string) string {
	for _, param := range payload.Parameters {
		if param.Name == name {
			return param.Value
		}
	}
	return ""
}

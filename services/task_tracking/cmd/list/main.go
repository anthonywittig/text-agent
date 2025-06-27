package main

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/ttacon/libphonenumber"
)

type Request struct {
	PhoneNumbers []string `json:"phone_numbers"`
}

type Response struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Tasks   []*task_repository.Task `json:"tasks"`
}

const (
	tableName = "text-agent-task-tracking"
)

func handleRequest(ctx context.Context, payload Request) (Response, error) {
	logger := zerolog.Ctx(ctx)

	repo, err := task_repository.New(ctx, tableName)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create repository")
		return Response{
			Status:  "internal_error",
			Message: "Internal error",
		}, nil
	}

	e164PhoneNumbers := make([]string, len(payload.PhoneNumbers))
	for i, phoneNumber := range payload.PhoneNumbers {
		number, err := libphonenumber.Parse(phoneNumber, "US")
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse phone number")
			return Response{
				Status:  "invalid_request",
				Message: "Unable to parse phone number " + phoneNumber,
			}, nil
		}
		e164PhoneNumbers[i] = libphonenumber.Format(number, libphonenumber.E164)
	}

	sort.Strings(e164PhoneNumbers)
	conversationID := strings.Join(e164PhoneNumbers, "_")

	logger.Info().Str("conversation_id", conversationID).Msg("Processing conversation")

	tasks, err := repo.ListTasksByConversation(conversationID)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to list tasks")
		return Response{
			Status:  "internal_error",
			Message: "Internal error",
		}, nil
	}

	return Response{
		Status:  "success",
		Message: "Tasks listed successfully",
		Tasks:   tasks,
	}, nil
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	requestWrapper := func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		requestID := "unknown"
		if lc != nil {
			requestID = lc.AwsRequestID
		}
		logger := logger.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Interface("request", payload).Msg("received request")

		var request Request
		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal request")
			response := Response{
				Status:  "invalid_request",
				Message: "Invalid request",
			}
			responseJSON, _ := json.Marshal(response)
			return responseJSON, nil
		}

		response, err := handleRequest(ctx, request)
		if err != nil {
			logger.Error().Err(err).Msg("failed to handle request")
			return nil, err
		}

		responseJSON, _ := json.Marshal(response)
		return responseJSON, nil
	}

	lambda.Start(requestWrapper)
}

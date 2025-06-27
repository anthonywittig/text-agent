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
	"github.com/rs/zerolog/log"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
	"github.com/ttacon/libphonenumber"
)

// Request represents the incoming API request
type Request struct {
	Data json.RawMessage `json:"data"`
}

// Response represents the API response
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

const (
	tableName = "text-agent-task-tracking"
)

func handleRequest(ctx context.Context, request Request) (Response, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := "unknown"
	if lc != nil {
		requestID = lc.AwsRequestID
	}
	logger := log.With().Str("request_id", requestID).Logger()
	ctx = logger.WithContext(ctx)

	logger.Info().Interface("request", request).Msg("Received request")

	repo, err := task_repository.New(ctx, tableName)
	if err != nil {
		return Response{Status: "error", Message: "Failed to create repository: " + err.Error()}, nil
	}

	var data struct {
		PhoneNumbers []string `json:"phone_numbers"`
	}
	if err := json.Unmarshal(request.Data, &data); err != nil {
		return Response{Status: "error", Message: "Invalid request data: " + err.Error()}, nil
	}

	e164PhoneNumbers := make([]string, len(data.PhoneNumbers))
	for i, phoneNumber := range data.PhoneNumbers {
		number, err := libphonenumber.Parse(phoneNumber, "US")
		if err != nil {
			return Response{Status: "error", Message: "Failed to parse phone number: " + err.Error()}, nil
		}
		e164PhoneNumbers[i] = libphonenumber.Format(number, libphonenumber.E164)
	}

	sort.Strings(e164PhoneNumbers)
	conversationID := strings.Join(e164PhoneNumbers, "_")

	log.Info().Str("conversation_id", conversationID).Msg("Processing conversation")

	tasks, err := repo.ListTasksByConversation(conversationID)
	if err != nil {
		log.Error().Err(err).Str("conversation_id", conversationID).Msg("Failed to list tasks")
		return Response{Status: "error", Message: "Failed to list items: " + err.Error()}, nil
	}

	tasksJSON, _ := json.Marshal(tasks)
	return Response{
		Status: "success",
		Data:   tasksJSON,
	}, nil
}

func main() {
	// Configure zerolog for Lambda environment
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	lambda.Start(handleRequest)
}

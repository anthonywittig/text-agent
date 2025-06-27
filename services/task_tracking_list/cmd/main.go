package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	log "github.com/pion/ion-log"

	"github.com/anthonywittig/text-agent/services/task_tracking/pkg/task_repository"
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
	// Create repository
	repo, err := task_repository.New(ctx, tableName)
	if err != nil {
		return Response{Status: "error", Message: "Failed to create repository: " + err.Error()}, nil
	}

	// Get the Lambda context to access request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := "unknown"
	if lc != nil {
		requestID = lc.AwsRequestID
	}

	// Log the incoming request with structured logging
	log.Infof("Received request: %s", requestID)

	var data struct {
		ConversationID string `json:"conversation_id"`
	}
	if err := json.Unmarshal(request.Data, &data); err != nil {
		return Response{Status: "error", Message: "Invalid request data: " + err.Error()}, nil
	}

	tasks, err := repo.ListTasksByConversation(data.ConversationID)
	if err != nil {
		return Response{Status: "error", Message: "Failed to list items: " + err.Error()}, nil
	}

	tasksJSON, _ := json.Marshal(tasks)
	return Response{
		Status: "success",
		Data:   tasksJSON,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	log "github.com/pion/ion-log"
)

func init() {
	log.Init("debug")
}

type Request struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func handleRequest(ctx context.Context, request Request) (Response, error) {
	// Get the Lambda context to access request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := "unknown"
	if lc != nil {
		requestID = lc.AwsRequestID
	}

	// Log the incoming request with structured logging
	log.Infof("Received request",
		"request_id", requestID,
		"action", request.Action,
		"data", string(request.Data),
	)

	// For now, just return a success response
	return Response{
		Status:  "success",
		Message: "Request received",
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/pion/ion-log"
	"github.com/google/uuid"
)

func init() {
	log.Init("debug")
	log.Infof("Initializing task tracking service")
}

// Task represents a single task in the tracking system
type Task struct {
	ID            string     `json:"id" dynamodbav:"id"`
	ConversationID string     `json:"conversation_id" dynamodbav:"conversation_id"`
	Name          string     `json:"name" dynamodbav:"name"`
	Description   string     `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Source        string     `json:"source" dynamodbav:"source"`
	Status        string     `json:"status" dynamodbav:"status"`
	DueDate       *time.Time `json:"due_date,omitempty" dynamodbav:"due_date,omitempty,unixtime"`
}

// Request represents the incoming API request
type Request struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// Response represents the API response
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

const (
	tableName = "TaskTracking"
)

// TaskManager handles task operations
type TaskManager struct {
	db *dynamodb.Client
}

// NewTaskManager creates a new TaskManager instance
func NewTaskManager() (*TaskManager, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create DynamoDB client
	db := dynamodb.NewFromConfig(cfg)

	return &TaskManager{
		db: db,
	}, nil
}

// CreateTask creates a new task
func (tm *TaskManager) CreateTask(conversationID, name, description, source string, dueDate *time.Time) (*Task, error) {
	task := Task{
		ID:            uuid.New().String(),
		ConversationID: conversationID,
		Name:          name,
		Description:   description,
		Source:        source,
		Status:        "open",
		DueDate:       dueDate,
	}


	// Convert task to DynamoDB attribute value map
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %v", err)
	}

	// Put item in DynamoDB
	_, err = tm.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return &task, nil
}

// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(id string) (*Task, error) {
	result, err := tm.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}
	if result.Item == nil {
		return nil, errors.New("task not found")
	}

	var task Task
	err = attributevalue.UnmarshalMap(result.Item, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %v", err)
	}

	return &task, nil
}

// UpdateTaskStatus updates a task's status
// Status should be one of: "open", "canceled", or "completed" (which will be formatted as "completed on {DATE}")
func (tm *TaskManager) UpdateTaskStatus(id, status string, completionDate *time.Time) (*Task, error) {
	// First get the task to ensure it exists
	task, err := tm.GetTask(id)
	if err != nil {
		return nil, err
	}

	// Validate status
	switch status {
	case "open", "canceled":
		task.Status = status
	case "completed":
		// If completion date is not provided, use current time
		completeTime := time.Now()
		if completionDate != nil {
			completeTime = *completionDate
		}
		// Format as "completed on {DATE}"
		task.Status = fmt.Sprintf("completed on %s", completeTime.Format("2006-01-02"))
	default:
		return nil, fmt.Errorf("invalid status: %s, must be 'open', 'canceled', or 'completed'", status)
	}

	// Convert task to DynamoDB attribute value map
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %v", err)
	}

	// Update item in DynamoDB
	_, err = tm.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return task, nil
}

// ListTasksByConversation returns all tasks for a specific conversation
func (tm *TaskManager) ListTasksByConversation(conversationID string) ([]Task, error) {
	// Note: This assumes a GSI on conversation_id is created in DynamoDB
	// In a production environment, you would use a query with the GSI
	// For now, we'll scan the table with a filter (not recommended for production)
	result, err := tm.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("conversation_id = :conversation_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":conversation_id": &types.AttributeValueMemberS{Value: conversationID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan DynamoDB: %v", err)
	}

	tasks := make([]Task, 0, len(result.Items))
	for _, item := range result.Items {
		var task Task
		err := attributevalue.UnmarshalMap(item, &task)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// DeleteTask removes a task by ID
func (tm *TaskManager) DeleteTask(id string) error {
	_, err := tm.db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %v", err)
	}

	return nil
}

var taskManager *TaskManager

func init() {
	var err error
	taskManager, err = NewTaskManager()
	if err != nil {
		log.Errorf("Failed to initialize task manager: %v", err)
	}
}

func handleRequest(ctx context.Context, request Request) (Response, error) {
	// Get the Lambda context to access request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := "unknown"
	if lc != nil {
		requestID = lc.AwsRequestID
	}

	// Log the incoming request with structured logging
	log.Infof("Received request: %s - %s", requestID, request.Action)


	if taskManager == nil {
		return Response{
			Status:  "error",
			Message: "Task manager not initialized",
		}, nil
	}

	switch request.Action {
	case "create_item":
		var data struct {
			ConversationID string     `json:"conversation_id"`
			Name          string     `json:"name"`
			Description   string     `json:"description,omitempty"`
			Source        string     `json:"source"`
			DueDate       *time.Time `json:"due_date,omitempty"`
		}
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return Response{Status: "error", Message: "Invalid request data: " + err.Error()}, nil
		}

		task, err := taskManager.CreateTask(
			data.ConversationID,
			data.Name,
			data.Description,
			data.Source,
			data.DueDate,
		)
		if err != nil {
			return Response{Status: "error", Message: "Failed to create item: " + err.Error()}, nil
		}

		taskJSON, _ := json.Marshal(task)
		return Response{
			Status:  "success",
			Message: "Item created successfully",
			Data:    taskJSON,
		}, nil

	case "list_items":
		var data struct {
			ConversationID string `json:"conversation_id"`
		}
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return Response{Status: "error", Message: "Invalid request data: " + err.Error()}, nil
		}

		tasks, err := taskManager.ListTasksByConversation(data.ConversationID)
		if err != nil {
			return Response{Status: "error", Message: "Failed to list items: " + err.Error()}, nil
		}

		tasksJSON, _ := json.Marshal(tasks)
		return Response{
			Status: "success",
			Data:   tasksJSON,
		}, nil

	case "delete_item":
		var data struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(request.Data, &data); err != nil {
			return Response{Status: "error", Message: "Invalid request data: " + err.Error()}, nil
		}

		if err := taskManager.DeleteTask(data.ID); err != nil {
			return Response{Status: "error", Message: "Failed to delete item: " + err.Error()}, nil
		}

		return Response{
			Status:  "success",
			Message: "Item deleted successfully",
		}, nil

	default:
		return Response{
			Status:  "error",
			Message: "Unknown action: " + request.Action,
		}, nil
	}
}

func main() {
	lambda.Start(handleRequest)
}

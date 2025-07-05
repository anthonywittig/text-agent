package task_repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// DynamoRepository manages task operations against DynamoDB
type DynamoRepository struct {
	db        *dynamodb.Client
	tableName string
}

// New creates a new DynamoRepository
func New(ctx context.Context, tableName string) (TaskRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %v", err)
	}

	// Create DynamoDB client
	db := dynamodb.NewFromConfig(cfg)
	return &DynamoRepository{
		db:        db,
		tableName: tableName,
	}, nil
}

// CreateTask creates a new task in DynamoDB
func (r *DynamoRepository) CreateTask(conversationId, name, description, source string) (*Task, error) {
	task := &Task{
		Id:             uuid.NewString(),
		ConversationId: conversationId,
		Name:           name,
		Description:    description,
		Source:         source,
	}

	av, err := attributevalue.MarshalMap(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	_, err = r.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put item in DynamoDB: %w", err)
	}

	task, err = r.GetTask(task.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task from DynamoDB: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (r *DynamoRepository) GetTask(id string) (*Task, error) {
	result, err := r.db.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("task not found with ID: %s", id)
	}

	var task Task
	err = attributevalue.UnmarshalMap(result.Item, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// DeleteTask removes a task by ID
func (r *DynamoRepository) DeleteTask(id string) error {
	_, err := r.db.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %w", err)
	}

	return nil
}

// ListTasksByConversation retrieves all tasks for a conversation
func (r *DynamoRepository) ListTasksByConversation(conversationID string) ([]*Task, error) {
	result, err := r.db.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("ConversationIdIndex"),
		KeyConditionExpression: aws.String("conversation_id = :convId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":convId": &types.AttributeValueMemberS{Value: conversationID},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query items from DynamoDB: %w", err)
	}

	if len(result.Items) == 0 {
		return []*Task{}, nil
	}

	var tasks []*Task
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}

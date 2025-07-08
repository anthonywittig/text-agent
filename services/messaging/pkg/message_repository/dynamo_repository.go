package message_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoRepository struct {
	db        *dynamodb.Client
	tableName string
}

func New(ctx context.Context, tableName string) (MessageRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %v", err)
	}

	db := dynamodb.NewFromConfig(cfg)
	return &DynamoRepository{
		db:        db,
		tableName: tableName,
	}, nil
}

func (r *DynamoRepository) CreateMessage(conversationId, from, body string) (*Message, error) {
	message := &Message{
		Id:             uuid.NewString(),
		ConversationId: conversationId,
		From:           from,
		Body:           body,
		SentAt:         time.Now().UnixMilli(),
	}

	av, err := attributevalue.MarshalMap(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = r.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put item in DynamoDB: %w", err)
	}

	message, err = r.GetMessage(message.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get message from DynamoDB: %w", err)
	}

	return message, nil
}

func (r *DynamoRepository) GetMessage(id string) (*Message, error) {
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

	var message Message
	err = attributevalue.UnmarshalMap(result.Item, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &message, nil
}

func (r *DynamoRepository) ListRecentMessagesByConversation(conversationID string) ([]*Message, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("ConversationIdIndex"),
		KeyConditionExpression: aws.String("conversation_id = :convId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":convId": &types.AttributeValueMemberS{Value: conversationID},
		},
		ScanIndexForward: aws.Bool(false), // true for ascending (oldest first), false for descending (newest first)
		Limit:            aws.Int32(20),
	}

	result, err := r.db.Query(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items from DynamoDB: %w", err)
	}

	if len(result.Items) == 0 {
		return []*Message{}, nil
	}

	var messages []*Message
	err = attributevalue.UnmarshalListOfMaps(result.Items, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	// Reverse the messages so that the oldest message comes first.
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

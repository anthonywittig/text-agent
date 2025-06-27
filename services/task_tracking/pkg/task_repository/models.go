package task_repository

import (
	"time"
)

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

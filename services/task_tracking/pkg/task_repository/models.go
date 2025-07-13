package task_repository

// Task represents a single task in the tracking system
type Task struct {
	Id             string `json:"id" dynamodbav:"id"`
	ConversationId string `json:"conversation_id" dynamodbav:"conversation_id"`
	Name           string `json:"name" dynamodbav:"name"`
	Description    string `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Source         string `json:"source" dynamodbav:"source"`
}

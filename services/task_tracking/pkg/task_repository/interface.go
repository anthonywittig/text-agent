package task_repository

import (
	"time"
)

// TaskRepository defines the interface for task storage operations
type TaskRepository interface {
	// CreateTask creates a new task
	CreateTask(conversationID, name, description, source string, dueDate *time.Time) (*Task, error)

	// UpdateTaskStatus updates the status of a task
	UpdateTaskStatus(id, status string, completionDate *time.Time) (*Task, error)

	// DeleteTask removes a task by ID
	DeleteTask(id string) error

	// ListTasksByConversation retrieves all tasks for a conversation
	ListTasksByConversation(conversationID string) ([]*Task, error)
}

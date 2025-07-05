package task_repository

// TaskRepository defines the interface for task storage operations
type TaskRepository interface {
	// CreateTask creates a new task
	CreateTask(conversationId, name, description, source string) (*Task, error)

	// DeleteTask removes a task by ID
	DeleteTask(id string) error

	// ListTasksByConversation retrieves all tasks for a conversation
	ListTasksByConversation(conversationID string) ([]*Task, error)
}

package message_repository

type MessageRepository interface {
	CreateMessage(conversationId, from, body string) (*Message, error)
	ListRecentMessagesByConversation(conversationID string) ([]*Message, error)
}

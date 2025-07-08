package message_repository

type Message struct {
	Id             string `json:"id" dynamodbav:"id"`
	ConversationId string `json:"conversation_id" dynamodbav:"conversation_id"`
	Body           string `json:"body" dynamodbav:"body"`
	From           string `json:"from" dynamodbav:"from"`
	SentAt         int64  `json:"sent_at" dynamodbav:"sent_at"` // UNIX timestamp in milliseconds
}

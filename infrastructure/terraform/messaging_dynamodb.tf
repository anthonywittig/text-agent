resource "aws_dynamodb_table" "messaging" {
  name         = "text-agent-messaging"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "conversation_id"
    type = "S"
  }

  attribute {
    name = "sent_at"
    type = "N"
  }

  global_secondary_index {
    name            = "ConversationIdIndex"
    hash_key        = "conversation_id"
    range_key       = "sent_at"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "SentAtIndex"
    hash_key        = "sent_at"
    projection_type = "ALL"
  }

  tags = {
    Name    = "text-agent-messaging"
    Service = "TextAgent"
  }
}

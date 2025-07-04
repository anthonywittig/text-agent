resource "aws_dynamodb_table" "task_tracking" {
  name         = "text-agent-task-tracking"
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

  global_secondary_index {
    name            = "ConversationIdIndex"
    hash_key        = "conversation_id"
    projection_type = "ALL"
  }
  
  tags = {
    Name    = "text-agent-task-tracking"
    Service = "TextAgent"
  }
}

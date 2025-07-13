resource "aws_bedrockagent_agent_action_group" "messaging" {
  agent_id      = aws_bedrockagent_agent.text_agent.agent_id
  agent_version = "DRAFT"

  action_group_name = "Messaging"

  function_schema {

    member_functions {
      functions {
        name        = "messaging_create"
        description = "Use this function to create a new message; use it when you need to send a message to the conversation. E.g. when you create a task or delete a task. Or when a user asks you a question in one of their messages."
        parameters {
          map_block_key = "conversation_phone_numbers"
          type          = "array"
          description   = "The phone numbers involved in the conversation"
          required      = true
        }
        parameters {
          map_block_key = "from"
          type          = "string"
          description   = "This is used to identify the sender of the message. Always set this value to 'Assistant'."
          required      = true
        }
        parameters {
          map_block_key = "body"
          type          = "string"
          description   = "The body of the message"
          required      = true
        }
      }

      functions {
        name        = "messaging_list_recent"
        description = "Use this function to get the list of recent messages for a conversation."
        parameters {
          map_block_key = "conversation_phone_numbers"
          type          = "array"
          description   = "The phone numbers involved in the conversation"
          required      = true
        }
      }
    }
  }

  action_group_executor {
    lambda = aws_lambda_function.messaging.arn
  }

  depends_on = [
    aws_bedrockagent_agent.text_agent,
    time_sleep.prepare_agent_sleep,
    aws_lambda_function.messaging
  ]
}

resource "aws_lambda_permission" "allow_bedrock_messaging" {
  statement_id  = "AllowBedrockToInvokeLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.messaging.function_name
  principal     = "bedrock.amazonaws.com"

  depends_on = [
    aws_lambda_function.messaging
  ]
}

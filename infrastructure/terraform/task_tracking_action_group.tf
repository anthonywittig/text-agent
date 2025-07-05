resource "aws_bedrockagent_agent_action_group" "task_tracking_create" {
  agent_id      = aws_bedrockagent_agent.text_agent.agent_id
  agent_version = "DRAFT"

  action_group_name = "TaskTrackingCreate"

  function_schema {
    member_functions {
      functions {
        name        = "task_tracking_create"
        description = "Task Tracking Create"
        parameters {
          map_block_key = "conversation_phone_numbers"
          type          = "array"
          description   = "The phone numbers involved in the conversation"
          required      = true
        }
        parameters {
          map_block_key = "name"
          type          = "string"
          description   = "The name of the task"
          required      = true
        }
        parameters {
          map_block_key = "description"
          type          = "string"
          description   = "The description of the task"
          required      = true
        }
        parameters {
          map_block_key = "source"
          type          = "string"
          description   = "The text of the message that triggered the task creation"
          required      = true
        }
      }
    }
  }

  action_group_executor {
    lambda = aws_lambda_function.task_tracking.arn
  }

  depends_on = [
    aws_bedrockagent_agent.text_agent,
    time_sleep.prepare_agent_sleep,
    aws_lambda_function.task_tracking
  ]
}

resource "aws_bedrockagent_agent_action_group" "task_tracking_delete" {
  agent_id      = aws_bedrockagent_agent.text_agent.agent_id
  agent_version = "DRAFT"

  action_group_name = "TaskTrackingDelete"

  function_schema {
    member_functions {
      functions {
        name        = "task_tracking_delete"
        description = "Task Tracking Delete"
        parameters {
          map_block_key = "task_id"
          type          = "string"
          description   = "The ID of the task to delete"
          required      = true
        }
      }
    }
  }

  action_group_executor {
    lambda = aws_lambda_function.task_tracking.arn
  }

  depends_on = [
    aws_bedrockagent_agent.text_agent,
    time_sleep.prepare_agent_sleep,
    aws_lambda_function.task_tracking
  ]
}

resource "aws_bedrockagent_agent_action_group" "task_tracking_list" {
  agent_id      = aws_bedrockagent_agent.text_agent.agent_id
  agent_version = "DRAFT"

  action_group_name = "TaskTrackingList"

  function_schema {
    member_functions {
      functions {
        name        = "task_tracking_list"
        description = "Task Tracking List"
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
    lambda = aws_lambda_function.task_tracking.arn
  }

  depends_on = [
    aws_bedrockagent_agent.text_agent,
    time_sleep.prepare_agent_sleep,
    aws_lambda_function.task_tracking
  ]
}

resource "aws_lambda_permission" "allow_bedrock" {
  statement_id  = "AllowBedrockToInvokeLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.task_tracking.function_name
  principal     = "bedrock.amazonaws.com"

  depends_on = [
    aws_lambda_function.task_tracking
  ]
}

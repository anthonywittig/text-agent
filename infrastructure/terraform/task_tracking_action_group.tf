# Create an action group for the Bedrock agent to enable task tracking tools
resource "aws_bedrockagent_agent_action_group" "task_tracking" {
  # Reference the agent ID and set the version to DRAFT
  agent_id      = aws_bedrockagent_agent.text_agent.agent_id
  agent_version = "DRAFT"

  # Define the action group
  action_group_name = "TaskTracking"

  # Inline schema that defines the available actions
  function_schema {
    member_functions {
      functions {
        name        = "task_tracking"
        description = "Task Tracking"
        parameters {
          map_block_key = "action"
          type          = "string"
          description   = "The action to perform"
          required      = true
        }
      }
    }
  }

  # Set up the Lambda function executor
  action_group_executor {
    lambda = aws_lambda_function.task_tracking_list.arn
  }

  # Ensure the agent and Lambda are created first
  depends_on = [
    aws_bedrockagent_agent.text_agent,
    time_sleep.prepare_agent_sleep,
    aws_lambda_function.task_tracking_list
  ]
}

# Allow Bedrock to invoke the Lambda function
resource "aws_lambda_permission" "allow_bedrock" {
  statement_id  = "AllowBedrockToInvokeLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.task_tracking_list.function_name
  principal     = "bedrock.amazonaws.com"

  depends_on = [
    aws_lambda_function.task_tracking_list
  ]
}

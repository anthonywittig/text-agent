# You probably need to request access to the agent (I did so manually).

resource "aws_bedrockagent_agent" "text_agent" {
  agent_name              = "text-agent-v1"
  agent_resource_role_arn = aws_iam_role.agent_role.arn
  foundation_model        = "arn:aws:bedrock:us-west-2:648145505041:inference-profile/us.amazon.nova-premier-v1:0"

  instruction = <<-EOT
    You are an AI assistant monitoring a group text conversation. Your tasks are:
    1. Track/update todo items mentioned in the conversation

    When a message is received, determine if you should:
    - Execute one of your tools
    - Generate a response to send back to the conversation

    Use the task tracking tools to manage tasks:
    - Use create_item when someone mentions a new task or commitment
    - Use list_items to check existing tasks for a conversation
    - Use delete_item when a task is no longer relevant
  EOT
}

resource "null_resource" "prepare_agent" {
  triggers = {
    agent_state = sha256(jsonencode(aws_bedrockagent_agent.text_agent))
  }

  provisioner "local-exec" {
    command = "aws bedrock-agent prepare-agent --agent-id ${aws_bedrockagent_agent.text_agent.agent_id} --region us-west-2 --profile ${var.aws_profile}"
  }
}

# Optional: Add a small delay to ensure preparation completes
resource "time_sleep" "prepare_agent_sleep" {
  create_duration = "5s"

  lifecycle {
    replace_triggered_by = [null_resource.prepare_agent]
  }
}

resource "aws_iam_role" "agent_role" {
  name = "text-agent-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "bedrock.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "agent_policy" {
  name = "text-agent-policy"
  role = aws_iam_role.agent_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "bedrock:*"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = aws_lambda_function.task_tracking_list.arn
      }
    ]
  })
}


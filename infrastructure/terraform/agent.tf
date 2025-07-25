# You probably need to request access to the agent (I did so manually).
# I manually enabled "Model invocation logging" in the AWS console.

resource "aws_cloudwatch_log_group" "bedrock_logs" {
  name              = "/aws/bedrock/text-agent"
  retention_in_days = 30

  tags = {
    Name    = "text-agent-bedrock-logs"
    Service = "TextAgent"
  }
}

resource "aws_bedrockagent_agent" "text_agent" {
  agent_name              = "text-agent-v1"
  agent_resource_role_arn = aws_iam_role.agent_role.arn
  foundation_model        = "arn:aws:bedrock:us-west-2:648145505041:inference-profile/us.amazon.nova-premier-v1:0"

  instruction = <<-EOT
    You are an AI assistant monitoring a group text conversation. Your primary task is to track/update todo items mentioned in the conversation.

    You will be invoked every time a new message is received. You'll want to:
    - Get the list of recent messages for the conversation.
    - Get the list of tasks.
    - Compare the tasks to the conversation and create/delete tasks if needed.
    - Send the users a message if appropriate (e.g. if a task is created or deleted, or if a user asks you a question).
  EOT
}

resource "null_resource" "prepare_agent" {
  triggers = {
    agent_state = sha256(jsonencode(aws_bedrockagent_agent.text_agent))
    #task_tracking_action_group = sha256(jsonencode(aws_bedrockagent_agent_action_group.task_tracking))
    #messaging_action_group     = sha256(jsonencode(aws_bedrockagent_agent_action_group.messaging))
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
        Resource = [
          aws_lambda_function.task_tracking.arn,
          aws_lambda_function.messaging.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Resource = [
          aws_cloudwatch_log_group.bedrock_logs.arn,
          "${aws_cloudwatch_log_group.bedrock_logs.arn}:*"
        ]
      }
    ]
  })
}

# Additional policy for Bedrock logging
resource "aws_iam_role_policy" "bedrock_logging_policy" {
  name = "text-agent-bedrock-logging-policy"
  role = aws_iam_role.agent_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Resource = [
          "arn:aws:logs:us-west-2:*:log-group:/aws/bedrock/*",
          "arn:aws:logs:us-west-2:*:log-group:/aws/bedrock/*:*"
        ]
      }
    ]
  })
}

# Force agent alias to update on every deployment
resource "null_resource" "force_alias_update" {
  triggers = {
    git_sha = var.git_sha
  }
}

resource "aws_bedrockagent_agent_alias" "text_agent_alias" {
  agent_id         = aws_bedrockagent_agent.text_agent.agent_id
  agent_alias_name = "production"

  # Force recreation when agent, action groups, or deployment changes
  lifecycle {
    replace_triggered_by = [
      aws_bedrockagent_agent.text_agent,
      aws_bedrockagent_agent_action_group.task_tracking,
      aws_bedrockagent_agent_action_group.messaging,
      null_resource.force_alias_update
    ]
  }
}

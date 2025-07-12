# CloudWatch Log Group with retention period
resource "aws_cloudwatch_log_group" "messaging" {
  name              = "/aws/lambda/text-agent-messaging"
  retention_in_days = 14
}

resource "aws_lambda_function" "messaging" {
  function_name = "text-agent-messaging"
  role          = aws_iam_role.lambda_exec_messaging.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.text_agent_messaging.repository_url}:${var.git_sha}"
  memory_size   = 128
  timeout       = 300
  architectures = ["arm64"]

  environment {
    variables = {
      AGENT_ALIAS_ID_SECRET_ID = aws_secretsmanager_secret.bedrock_agent_alias_id.id
      AGENT_ID_SECRET_ID       = aws_secretsmanager_secret.bedrock_agent_id.id
    }
  }

  depends_on = [
    aws_iam_role_policy.lambda_exec_policy_messaging,
    aws_cloudwatch_log_group.messaging,
  ]
}

resource "aws_iam_role" "lambda_exec_messaging" {
  name = "text-agent-messaging-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_exec_policy_messaging" {
  name = "text-agent-messaging-exec-policy"
  role = aws_iam_role.lambda_exec_messaging.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "bedrock:InvokeAgent"
        ]
        Resource = [
          "*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "xray:PutTraceSegments",
          "xray:PutTelemetryRecords"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:*",
        ]
        Resource = [
          aws_dynamodb_table.messaging.arn,
          "${aws_dynamodb_table.messaging.arn}/index/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          aws_secretsmanager_secret.bedrock_agent_alias_id.arn,
          aws_secretsmanager_secret.bedrock_agent_id.arn,
        ]
      }
    ]
  })
}

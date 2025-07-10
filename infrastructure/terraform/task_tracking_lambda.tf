# CloudWatch Log Group with retention period
resource "aws_cloudwatch_log_group" "task_tracking" {
  name              = "/aws/lambda/text-agent-task-tracking"
  retention_in_days = 14
}

resource "aws_lambda_function" "task_tracking" {
  function_name = "text-agent-task-tracking"
  role          = aws_iam_role.lambda_exec.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.text_agent_task_tracking.repository_url}:${var.git_sha}"
  memory_size   = 128
  timeout       = 30
  architectures = ["arm64"]

  environment {
    variables = {}
  }

  depends_on = [
    aws_iam_role_policy.lambda_exec_policy,
    aws_cloudwatch_log_group.task_tracking,
  ]
}

resource "aws_iam_role" "lambda_exec" {
  name = "text-agent-task-tracking-exec-role"

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

resource "aws_iam_role_policy" "lambda_exec_policy" {
  name = "text-agent-task-tracking-exec-policy"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
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
          aws_dynamodb_table.task_tracking.arn,
          "${aws_dynamodb_table.task_tracking.arn}/index/*"
        ]
      }
    ]
  })
}

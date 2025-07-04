# CloudWatch Log Group with retention period
resource "aws_cloudwatch_log_group" "task_tracking_list" {
  name              = "/aws/lambda/text-agent-task-tracking-list"
  retention_in_days = 14
}

resource "aws_lambda_function" "task_tracking_list" {
  function_name = "text-agent-task-tracking-list"
  role          = aws_iam_role.lambda_exec.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.text_agent_task_tracking_list.repository_url}:${var.git_sha}"
  memory_size   = 128
  timeout       = 30
  architectures = ["arm64"]

  environment {
    variables = {}
  }

  depends_on = [
    aws_iam_role_policy.lambda_exec_policy,
    aws_cloudwatch_log_group.task_tracking_list,
  ]
}

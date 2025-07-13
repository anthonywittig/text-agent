resource "aws_ecr_repository" "text_agent_messaging" {
  name = "text-agent-messaging"
}

resource "aws_ecr_lifecycle_policy" "text_agent_messaging" {
  repository = aws_ecr_repository.text_agent_messaging.name
  policy     = <<EOF
  {
    "rules": [
      {
        "rulePriority": 1,
        "description": "Expire older images.",
        "selection": {
          "tagStatus": "any",
          "countType": "imageCountMoreThan",
          "countNumber": 1
        },
        "action": {
          "type": "expire"
        }
      }
    ]
  }
EOF
}

resource "aws_ecr_repository" "text_agent_task_tracking" {
  name = "text-agent-task-tracking"
}

resource "aws_ecr_lifecycle_policy" "text_agent_task_tracking" {
  repository = aws_ecr_repository.text_agent_task_tracking.name
  policy     = <<EOF
  {
    "rules": [
      {
        "rulePriority": 1,
        "description": "Expire older images.",
        "selection": {
          "tagStatus": "any",
          "countType": "imageCountMoreThan",
          "countNumber": 1
        },
        "action": {
          "type": "expire"
        }
      }
    ]
  }
EOF
}

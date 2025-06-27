#!/bin/bash

# Exit on error
set -euo pipefail

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT"

export "$(grep -v '^#' ./.env | xargs)"

# Check if AWS_PROFILE is set
if [ -z "$AWS_PROFILE" ]; then
    echo "Error: AWS_PROFILE environment variable is not set"
    exit 1
fi

# Get the AWS account ID and region
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)

# Get the current git commit hash
GIT_COMMIT=$(git rev-parse --short HEAD)

# Get the repository name from the ECR resource
REPO_NAME="text-agent-task-tracking-list"

# Get the full ECR repository URL
ECR_REPO="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPO_NAME}"

# Login to ECR
aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${ECR_REPO}"

cd services/task_tracking_list

# Build the Docker image
DOCKER_BUILDKIT=1 docker build \
  -t "${ECR_REPO}":"${GIT_COMMIT}" \
  -t "${ECR_REPO}":latest \
  -f cmd/Dockerfile \
  .

# Push the image to ECR
docker push "${ECR_REPO}":"${GIT_COMMIT}"
docker push "${ECR_REPO}":latest

# Print success message
echo "Successfully built and pushed Docker image to ${ECR_REPO}:${GIT_COMMIT} and ${ECR_REPO}:latest"

#!/bin/bash

# Exit on error
set -euo pipefail

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT"

export "$(grep -v '^#' ./.env | xargs)"

if [ -z "$AWS_PROFILE" ]; then
    echo "Error: AWS_PROFILE environment variable is not set"
    exit 1
fi

# Get the AWS account ID and region
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)

# Get the current git commit hash
GIT_COMMIT=$(git rev-parse --short HEAD)

###
# Messaging
###

cd services/messaging
REPO_NAME="text-agent-messaging"
ECR_REPO="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPO_NAME}"
aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${ECR_REPO}"
DOCKER_BUILDKIT=1 docker build \
  -t "${ECR_REPO}":"${GIT_COMMIT}" \
  -t "${ECR_REPO}":latest \
  -f cmd/Dockerfile \
  .
docker push "${ECR_REPO}":"${GIT_COMMIT}"
docker push "${ECR_REPO}":latest
echo "Successfully built and pushed Docker image to ${ECR_REPO}:${GIT_COMMIT} and ${ECR_REPO}:latest"
cd ../../

###
# Task Tracking
###

cd services/task_tracking
REPO_NAME="text-agent-task-tracking"
ECR_REPO="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPO_NAME}"
aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${ECR_REPO}"
DOCKER_BUILDKIT=1 docker build \
  -t "${ECR_REPO}":"${GIT_COMMIT}" \
  -t "${ECR_REPO}":latest \
  -f cmd/Dockerfile \
  .
docker push "${ECR_REPO}":"${GIT_COMMIT}"
docker push "${ECR_REPO}":latest
echo "Successfully built and pushed Docker image to ${ECR_REPO}:${GIT_COMMIT} and ${ECR_REPO}:latest"
cd ../../

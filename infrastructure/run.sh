#!/bin/bash

set -euo pipefail

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT"/infrastructure/terraform

export "$(grep -v '^#' ../../.env | xargs)"

export TF_VAR_git_sha=$(git rev-parse --short HEAD)
export TF_VAR_aws_profile="$AWS_PROFILE"

terraform init -backend-config=backend.hcl
terraform plan -out=plan.out

read -p "Are you sure you want to apply the changes? [y/N]: " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
  echo "Exiting"
  exit 0
fi

terraform apply plan.out

rm plan.out

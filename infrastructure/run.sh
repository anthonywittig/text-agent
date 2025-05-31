#!/bin/bash

set -euo pipefail

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT"/infrastructure/terraform

export "$(grep -v '^#' ../.env_init | xargs)"

terraform init -backend-config=backend.hcl
export "$(grep -v '^#' ../.env_plan | xargs)"
terraform plan -out=plan.out

read -p "Are you sure you want to apply the changes? (y/N): " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
  echo "Exiting"
  exit 0
fi

terraform apply plan.out

rm plan.out

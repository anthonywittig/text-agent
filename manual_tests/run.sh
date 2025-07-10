#! /bin/bash

export "$(grep -v '^#' ../.env | xargs)"

if [ -z "$AWS_PROFILE" ]; then
    echo "Error: AWS_PROFILE environment variable is not set"
    exit 1
fi

#cd ../infrastructure/terraform || exit
#export TF_VAR_git_sha=$(git rev-parse --short HEAD)
#export TF_VAR_aws_profile="$AWS_PROFILE"
#AGENT_ALIAS_ID=$(terraform output -raw agent_alias_id)
#AGENT_ID=$(terraform output -raw agent_id)
#REGION=$(terraform output -raw aws_region)
#cd ../../manual_tests || exit

./bootstrap.sh

source .venv/bin/activate
python project/test_01.py

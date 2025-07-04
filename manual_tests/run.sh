#! /bin/bash

export "$(grep -v '^#' ../.env | xargs)"

if [ -z "$AWS_PROFILE" ]; then
    echo "Error: AWS_PROFILE environment variable is not set"
    exit 1
fi

cd ../infrastructure/terraform || exit

export TF_VAR_git_sha=$(git rev-parse --short HEAD)
export TF_VAR_aws_profile="$AWS_PROFILE"

AGENT_ALIAS_ID=$(terraform output -raw agent_alias_id)
AGENT_ID=$(terraform output -raw agent_id)
REGION=$(terraform output -raw aws_region)

cd ../../manual_tests || exit

./bootstrap.sh

source .venv/bin/activate
python project/test_01.py "$AGENT_ALIAS_ID" "$AGENT_ID"

# The `sed` removes the extension from the phone number.
#PHONE_NUMBER_1=$(faker -l en_US phone_number | sed 's/x[0-9]*//')
#PHONE_NUMBER_2=$(faker -l en_US phone_number | sed 's/x[0-9]*//')

#PHONE_NUMBERS="$PHONE_NUMBER_1,$PHONE_NUMBER_2"
#echo "PHONE_NUMBERS: $PHONE_NUMBERS"

#aws bedrock-agent-runtime invoke-agent \
#  --agent-id "$AGENT_ID" \
#  --agent-alias-id "$AGENT_ALIAS_ID" \
#  --session-id "$SESSION_ID" \
#  --input-text "context: phone numbers $PHONE_NUMBERS; message: Joe, please buy a new laptop for the office" \
#  --region "$REGION" \
#  --enable-trace

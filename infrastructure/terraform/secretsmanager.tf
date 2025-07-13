resource "aws_secretsmanager_secret" "bedrock_agent_alias_id" {
  name = "text-agent-bedrock-agent-alias-id"
}

resource "aws_secretsmanager_secret_version" "bedrock_agent_alias_id" {
  secret_id     = aws_secretsmanager_secret.bedrock_agent_alias_id.id
  secret_string = aws_bedrockagent_agent_alias.text_agent_alias.agent_alias_id
}

resource "aws_secretsmanager_secret" "bedrock_agent_id" {
  name = "text-agent-bedrock-agent-id"
}

resource "aws_secretsmanager_secret_version" "bedrock_agent_id" {
  secret_id     = aws_secretsmanager_secret.bedrock_agent_id.id
  secret_string = aws_bedrockagent_agent.text_agent.agent_id
}

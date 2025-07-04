output "agent_id" {
  description = "The ID of the Bedrock agent"
  value       = aws_bedrockagent_agent.text_agent.agent_id
}

output "agent_alias_id" {
  description = "The ID of the Bedrock agent alias"
  value       = aws_bedrockagent_agent_alias.text_agent_alias.agent_alias_id
}

output "aws_region" {
  description = "The AWS region"
  value       = local.region
}

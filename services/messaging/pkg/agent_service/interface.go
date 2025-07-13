package agent_service

import "context"

type AgentService interface {
	InvokeAgent(ctx context.Context, input string) error
}

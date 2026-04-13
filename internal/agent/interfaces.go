package agent

type AgentPolicy interface {
	Decide(ctx *AgentContext) (string, map[string]interface{}, bool)
}

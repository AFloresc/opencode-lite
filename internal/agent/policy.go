package agent

type AgentPolicy interface {
	Decide(ctx *AgentContext) (toolName string, args map[string]interface{}, done bool)
}

type SimplePolicy struct{}

func (p SimplePolicy) Decide(ctx *AgentContext) (string, map[string]interface{}, bool) {
	if ctx.LastResult.Error != "" {
		return "", nil, true
	}

	if ctx.Goal == "listar archivos" {
		return "file_tree", map[string]interface{}{"root": "workspace"}, false
	}

	return "", nil, true
}

package agent

func (ctx *AgentContext) Remember(key string, value interface{}) {
	ctx.Memory[key] = value
}

func (ctx *AgentContext) Recall(key string) interface{} {
	return ctx.Memory[key]
}

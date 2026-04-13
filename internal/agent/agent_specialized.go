package agent

type SpecializedAgent interface {
	Name() string
	CanHandle(goal string) bool
	Run(goal string, ctx *AgentContext) AgentContext
}

type BaseSpecializedAgent struct {
	name    string
	runtime *AgentRuntime
	matchFn func(goal string) bool
}

func (a *BaseSpecializedAgent) Name() string { return a.name }

func (a *BaseSpecializedAgent) CanHandle(goal string) bool {
	if a.matchFn == nil {
		return false
	}
	return a.matchFn(goal)
}

func (a *BaseSpecializedAgent) Run(goal string, ctx *AgentContext) AgentContext {
	// Reutilizamos el runtime, pero podemos pasar contexto inicial si quieres
	return a.runtime.Run(goal)
}

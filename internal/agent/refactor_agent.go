package agent

import "strings"

type RefactorAgent struct {
	name    string
	runtime *AgentRuntime
}

func NewRefactorAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewRefactorPolicy()

	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem

	return &RefactorAgent{
		name:    "refactor",
		runtime: rt,
	}
}

func (a *RefactorAgent) Name() string {
	return a.name
}

func (a *RefactorAgent) CanHandle(goal string) bool {
	g := strings.ToLower(goal)

	return strings.Contains(g, "refactor") ||
		strings.Contains(g, "limpiar código") ||
		strings.Contains(g, "mejorar estructura") ||
		strings.Contains(g, "optimizar") ||
		strings.Contains(g, "simplificar")
}

func (a *RefactorAgent) Run(goal string, ctx *AgentContext) AgentContext {
	return a.runtime.Run(goal)
}

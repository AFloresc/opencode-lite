package agent

import "strings"

type DocsAgent struct {
	name    string
	runtime *AgentRuntime
}

func NewDocsAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewDocsPolicy()

	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem

	return &DocsAgent{
		name:    "docs",
		runtime: rt,
	}
}

func (a *DocsAgent) Name() string {
	return a.name
}

func (a *DocsAgent) CanHandle(goal string) bool {
	g := strings.ToLower(goal)

	return strings.Contains(g, "documentar") ||
		strings.Contains(g, "comentarios") ||
		strings.Contains(g, "explicar") ||
		strings.Contains(g, "resumen") ||
		strings.Contains(g, "documentación")
}

func (a *DocsAgent) Run(goal string, ctx *AgentContext) AgentContext {
	return a.runtime.Run(goal)
}

package agent

import "strings"

type AnalysisAgent struct {
	name    string
	runtime *AgentRuntime
}

func NewAnalysisAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewAnalysisPolicy()

	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem

	return &AnalysisAgent{
		name:    "analysis",
		runtime: rt,
	}
}

func (a *AnalysisAgent) Name() string {
	return a.name
}

func (a *AnalysisAgent) CanHandle(goal string) bool {
	g := strings.ToLower(goal)

	return strings.Contains(g, "analizar") ||
		strings.Contains(g, "entender") ||
		strings.Contains(g, "dependencias") ||
		strings.Contains(g, "métricas") ||
		strings.Contains(g, "estructura") ||
		strings.Contains(g, "diagnóstico")
}

// Run cumple la interfaz SpecializedAgent.
// De momento ignora ctx y delega en el runtime.
func (a *AnalysisAgent) Run(goal string, ctx *AgentContext) AgentContext {
	return a.runtime.Run(goal)
}

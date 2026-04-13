package agent

import "strings"

func NewAnalysisAgent(projectID string, llm LLMClient) SpecializedAgent {
	policy := NewAnalysisPolicy() // tu policy de análisis
	rt := NewAgentRuntime(projectID, policy, llm)

	return &BaseSpecializedAgent{
		name:    "analysis",
		runtime: rt,
		matchFn: func(goal string) bool {
			g := strings.ToLower(goal)
			return strings.Contains(g, "analizar") ||
				strings.Contains(g, "entender") ||
				strings.Contains(g, "dependencias") ||
				strings.Contains(g, "métricas")
		},
	}
}

func NewRefactorAgent(projectID string, llm LLMClient) SpecializedAgent {
	policy := NewRefactorPolicy() // otra policy
	rt := NewAgentRuntime(projectID, policy, llm)

	return &BaseSpecializedAgent{
		name:    "refactor",
		runtime: rt,
		matchFn: func(goal string) bool {
			g := strings.ToLower(goal)
			return strings.Contains(g, "refactor") ||
				strings.Contains(g, "limpiar código") ||
				strings.Contains(g, "mejorar estructura")
		},
	}
}

func NewDocsAgent(projectID string, llm LLMClient) SpecializedAgent {
	policy := NewDocsPolicy()
	rt := NewAgentRuntime(projectID, policy, llm)

	return &BaseSpecializedAgent{
		name:    "docs",
		runtime: rt,
		matchFn: func(goal string) bool {
			g := strings.ToLower(goal)
			return strings.Contains(g, "documentar") ||
				strings.Contains(g, "comentarios") ||
				strings.Contains(g, "explicar") ||
				strings.Contains(g, "resumen")
		},
	}
}

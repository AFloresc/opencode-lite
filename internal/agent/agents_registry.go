package agent

import "strings"

func NewAnalysisAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewAnalysisPolicy()

	// Runtime cognitivo
	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem // <-- integración de memoria cognitiva

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

func NewRefactorAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewRefactorPolicy() // policy específica de refactor

	// Runtime cognitivo
	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem // <-- integración de memoria cognitiva avanzada

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

func NewDocsAgent(projectID string, llm LLMClient, mem *CognitiveMemory) SpecializedAgent {
	policy := NewDocsPolicy()

	// Runtime cognitivo
	rt := NewAgentRuntime(projectID, policy, llm)
	rt.Memory = mem // <-- integración de memoria cognitiva avanzada

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

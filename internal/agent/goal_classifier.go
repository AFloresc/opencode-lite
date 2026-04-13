package agent

import "strings"

type GoalClassifier struct {
	llm LLMClient
}

func NewGoalClassifier(llm LLMClient) *GoalClassifier {
	return &GoalClassifier{llm: llm}
}

// Devuelve el nombre del agente recomendado: "analysis", "refactor", "docs", etc.
func (c *GoalClassifier) Classify(goal string) string {
	prompt := `
Eres un clasificador de tareas para un sistema multi-agente.
Tu trabajo es decidir qué tipo de agente debe manejar este objetivo.

Agentes disponibles:
- analysis: análisis, dependencias, métricas, comprensión del código
- refactor: mejorar código, limpiar, reorganizar, extraer funciones
- docs: documentar, explicar, generar comentarios, resúmenes

Devuelve SOLO el nombre del agente, sin texto adicional.

Objetivo: "` + goal + `"
`

	resp, err := c.llm.Complete(prompt)
	if err != nil {
		// fallback simple
		g := strings.ToLower(goal)
		switch {
		case strings.Contains(g, "refactor"):
			return "refactor"
		case strings.Contains(g, "document"):
			return "docs"
		default:
			return "analysis"
		}
	}

	out := strings.ToLower(strings.TrimSpace(resp))
	if out == "analysis" || out == "refactor" || out == "docs" {
		return out
	}

	// fallback si el LLM devuelve algo raro
	return "analysis"
}

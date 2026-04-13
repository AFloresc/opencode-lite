package agent

import "strings"

func isTooBroad(goal string) bool {
	g := strings.ToLower(goal)
	return strings.Contains(g, "mejora el proyecto") ||
		strings.Contains(g, "arregla todo") ||
		strings.Contains(g, "optimiza todo") ||
		len(strings.Split(goal, " ")) < 3
}

func splitGoal(goal string) []string {
	// División simple por ahora
	return []string{
		"analizar estructura del proyecto",
		"detectar problemas principales",
		"proponer mejoras",
	}
}

func isGoalSatisfied(goal string, ctx *AgentContext) bool {
	// Heurística simple: si la última tool no devolvió nada útil → terminado
	if ctx.LastResult.Result == nil {
		return true
	}
	return false
}

func classifyGoal(goal string, llm LLMClient) string {
	prompt := `
Clasifica este objetivo en uno de estos agentes:
- analysis
- refactor
- docs

Devuelve solo el nombre del agente.

Objetivo: "` + goal + `"
`
	resp, err := llm.Complete(prompt)
	if err != nil {
		return "analysis"
	}

	resp = strings.ToLower(strings.TrimSpace(resp))
	if resp == "analysis" || resp == "refactor" || resp == "docs" {
		return resp
	}

	return "analysis"
}

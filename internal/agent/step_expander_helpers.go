package agent

import (
	"encoding/json"
	"opencode-lite/internal/tools"
	"strings"
)

// ------------------------------------------------------------
// 1. extractFilesFromResult
// ------------------------------------------------------------
// Extrae []string desde cualquier formato que devuelva una tool.
// Compatible con:
// - []string
// - []interface{}
// - map[string]interface{}{"tree": []string}
// - map[string]interface{}{"tree": []interface{}}
// ------------------------------------------------------------
func extractFilesFromResult(result tools.ToolResult) []string {
	// Caso 1: Result ya es []string
	if files, ok := result.Result.([]string); ok {
		return files
	}

	// Caso 2: Result es []interface{}
	if arr, ok := result.Result.([]interface{}); ok {
		out := make([]string, 0, len(arr))
		for _, v := range arr {
			if s, ok := v.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}

	// Caso 3: Result es map con "tree"
	if m, ok := result.Result.(map[string]interface{}); ok {
		if raw, ok := m["tree"]; ok {

			// tree: []string
			if files, ok := raw.([]string); ok {
				return files
			}

			// tree: []interface{}
			if arr, ok := raw.([]interface{}); ok {
				out := make([]string, 0, len(arr))
				for _, v := range arr {
					if s, ok := v.(string); ok {
						out = append(out, s)
					}
				}
				return out
			}
		}
	}

	return nil
}

// ------------------------------------------------------------
// 2. hasCycles
// ------------------------------------------------------------
// Heurística simple: detecta ciclos en dependencias buscando
// palabras clave en el resultado serializado.
// ------------------------------------------------------------
func hasCycles(result tools.ToolResult) bool {
	b, err := json.Marshal(result.Result)
	if err != nil {
		return false
	}
	s := strings.ToLower(string(b))
	return strings.Contains(s, "cycle") ||
		strings.Contains(s, "ciclo") ||
		strings.Contains(s, "dependency cycle")
}

// ------------------------------------------------------------
// 3. hasLongFunctions
// ------------------------------------------------------------
// Detecta funciones largas según el resultado de la tool.
// Compatible con herramientas que devuelven:
// - lista de funciones con tamaños
// - flags como "long_function"
// ------------------------------------------------------------
func hasLongFunctions(result tools.ToolResult) bool {
	b, err := json.Marshal(result.Result)
	if err != nil {
		return false
	}
	s := strings.ToLower(string(b))

	// Heurísticas comunes
	return strings.Contains(s, "long_function") ||
		strings.Contains(s, "large_function") ||
		strings.Contains(s, "func_too_long") ||
		strings.Contains(s, "over_200_lines")
}

// ------------------------------------------------------------
// 4. callLLMForSubplan
// ------------------------------------------------------------
// Permite que el LLM refine el plan dinámicamente.
// Usa el LLMClient si está disponible en ctx.Memory.
// ------------------------------------------------------------
func callLLMForSubplan(step PlanStep, result tools.ToolResult, ctx *AgentContext) []PlanStep {
	llm, ok := ctx.Memory["llm_client"].(LLMClient)
	if !ok || llm == nil {
		return nil
	}

	// Construimos un goal refinado
	b, _ := json.Marshal(result.Result)
	goal := "Refina el plan para el paso \"" + step.Description + "\" usando este contexto: " + string(b)

	steps, err := llm.ProposePlan(goal)
	if err != nil || len(steps) == 0 {
		return nil
	}

	out := make([]PlanStep, 0, len(steps))
	for _, s := range steps {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, PlanStep{Description: s})
		}
	}
	return out
}

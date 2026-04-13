package agent

import (
	"opencode-lite/internal/tools"
	"strings"
)

type StepExpander interface {
	Expand(step PlanStep, result tools.ToolResult, ctx *AgentContext) []PlanStep
}

type DefaultStepExpander struct{}

func NewDefaultStepExpander() *DefaultStepExpander {
	return &DefaultStepExpander{}
}

func (e *DefaultStepExpander) Expand(step PlanStep, result tools.ToolResult, ctx *AgentContext) []PlanStep {
	s := strings.ToLower(step.Description)

	// === 1. Si listar archivos devuelve muchos archivos → subplan ===
	if s == "listar archivos" {
		files := extractFilesFromResult(result)
		if len(files) > 500 {
			return []PlanStep{
				{"detectar hotspots"},
				{"buscar duplicación"},
				{"buscar funciones largas"},
			}
		}
	}

	// === 2. Si detectar dependencias encuentra ciclos → subplan ===
	if s == "detectar dependencias" {
		if hasCycles(result) {
			return []PlanStep{
				{"analizar ciclos de dependencias"},
				{"sugerir refactor de paquetes"},
			}
		}
	}

	// === 3. Si extraer funciones encuentra funciones enormes → subplan ===
	if s == "extraer funciones" {
		if hasLongFunctions(result) {
			return []PlanStep{
				{"buscar funciones largas"},
				{"sugerir refactor"},
			}
		}
	}

	// === 4. Si el LLM quiere refinar el plan ===
	if ctx.Memory["llm_refine"] == true {
		return callLLMForSubplan(step, result, ctx)
	}

	return nil
}

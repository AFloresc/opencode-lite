package agent

import (
	"strings"
)

//
// ============================================================
//  MetaEvaluation
// ============================================================
//

type MetaEvaluation struct {
	Confidence float64
	Flags      []string
	Advice     string
}

//
// ============================================================
//  Metacognition
// ============================================================
//

type Metacognition struct {
	llm LLMClient
}

func NewMetacognition(llm LLMClient) *Metacognition {
	return &Metacognition{llm: llm}
}

//
// ============================================================
//  Evaluate: núcleo de la metacognición
// ============================================================
//

func (m *Metacognition) Evaluate(goal string, rt *AgentRuntime, ctx *AgentContext) MetaEvaluation {
	flags := []string{}

	// ============================================================
	// 1. Señales del ExecutionMonitor
	// ============================================================

	if rt.Monitor.RepeatCount >= 2 {
		flags = append(flags, "loop_detected")
	}

	if rt.Monitor.FailureCount >= 2 {
		flags = append(flags, "repeated_failures")
	}

	if rt.Monitor.StallCount >= 2 {
		flags = append(flags, "stalled")
	}

	// ============================================================
	// 2. Goal ambiguo
	// ============================================================

	words := strings.Fields(goal)
	if len(words) < 3 {
		flags = append(flags, "ambiguous_goal")
	}

	// ============================================================
	// 3. Falta de progreso real
	// ============================================================

	if ctx.LastResult.Result == nil {
		flags = append(flags, "no_progress")
	}

	// ============================================================
	// 4. Patrones desde memoria cognitiva
	// ============================================================

	if rt.Memory.Recall("dependency_cycles") == true {
		flags = append(flags, "dependency_cycles_known")
	}

	if rt.Memory.Recall("long_functions_detected") == true {
		flags = append(flags, "long_functions_known")
	}

	if fails, ok := rt.Memory.Recall("fail_count").(int); ok && fails >= 3 {
		flags = append(flags, "high_failure_memory")
	}

	if success, ok := rt.Memory.Recall("success_count").(int); ok && success >= 5 {
		flags = append(flags, "high_success_memory")
	}

	// ============================================================
	// 5. Detectar modo incorrecto del planner o grounder
	// ============================================================

	if rt.Planner != nil {
		if hp, ok := rt.Planner.(*HybridPlanner); ok {
			if hp.Mode == "coarse" && rt.Monitor.FailureCount >= 2 {
				flags = append(flags, "planner_too_coarse")
			}
			if hp.Mode == "fine" && rt.Monitor.RepeatCount >= 2 {
				flags = append(flags, "planner_too_fine")
			}
		}
	}

	if rt.Grounder != nil {
		if cg, ok := rt.Grounder.(*ContextualToolGrounder); ok {
			if cg.Mode == "strict" && rt.Monitor.StallCount >= 2 {
				flags = append(flags, "grounder_too_strict")
			}
			if cg.Mode == "flexible" && rt.Monitor.FailureCount >= 2 {
				flags = append(flags, "grounder_too_flexible")
			}
		}
	}

	// ============================================================
	// 6. Autoevaluación LLM
	// ============================================================

	prompt := `
Eres un módulo de metacognición. Evalúa el rendimiento del agente.

Objetivo: "` + goal + `"
Flags detectadas: ` + strings.Join(flags, ", ") + `

Devuelve SOLO un JSON con:
{
  "confidence": número entre 0 y 1,
  "advice": "texto breve"
}
`

	resp, err := m.llm.Complete(prompt)
	if err != nil {
		return MetaEvaluation{
			Confidence: 0.5,
			Flags:      flags,
			Advice:     "Revisar estrategia actual.",
		}
	}

	return MetaEvaluation{
		Confidence: extractConfidence(resp),
		Flags:      flags,
		Advice:     extractAdvice(resp),
	}
}

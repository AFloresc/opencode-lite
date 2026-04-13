package agent

import "strings"

type MetaEvaluation struct {
	Confidence float64
	Flags      []string
	Advice     string
}

type Metacognition struct {
	llm LLMClient
}

func NewMetacognition(llm LLMClient) *Metacognition {
	return &Metacognition{llm: llm}
}

func (m *Metacognition) Evaluate(goal string, rt *AgentRuntime, ctx *AgentContext) MetaEvaluation {
	flags := []string{}

	// 1. Detectar loops
	if rt.Monitor.RepeatCount >= 2 {
		flags = append(flags, "loop_detected")
	}

	// 2. Detectar fallos repetidos
	if rt.Monitor.FailureCount >= 2 {
		flags = append(flags, "repeated_failures")
	}

	// 3. Detectar estancamiento
	if rt.Monitor.StallCount >= 2 {
		flags = append(flags, "stalled")
	}

	// 4. Detectar goal ambiguo
	if len(strings.Split(goal, " ")) < 3 {
		flags = append(flags, "ambiguous_goal")
	}

	// 5. Detectar falta de progreso
	if ctx.LastResult.Result == nil {
		flags = append(flags, "no_progress")
	}

	// 6. Detectar patrones desde memoria
	if rt.Memory.Recall("dependency_cycles") == true {
		flags = append(flags, "dependency_cycles_known")
	}

	if rt.Memory.Recall("long_functions_detected") == true {
		flags = append(flags, "long_functions_known")
	}

	// 7. LLM produce una auto‑evaluación
	prompt := `
Eres un módulo de metacognición. Evalúa el rendimiento del agente.

Objetivo: "` + goal + `"
Flags: ` + strings.Join(flags, ", ") + `

Devuelve:
- un nivel de confianza entre 0 y 1
- un consejo breve para mejorar la estrategia actual
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

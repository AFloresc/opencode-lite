package agent

import (
	"strings"
)

type Supervisor struct {
	llm           *LLMClient
	Metacognition *Metacognition
	Strategy      *StrategyEngine
	AOC           *AOC
}

func NewSupervisor(llm LLMClient, mem *CognitiveMemory) *Supervisor {
	return &Supervisor{
		llm:           &llm,
		Metacognition: NewMetacognition(llm),
		Strategy:      NewStrategyEngine(llm),
		AOC:           NewAOC(mem),
	}
}

type SupervisorDecision struct {
	Action    string   // "delegate", "clarify", "split", "finish", "replan"
	AgentName string   // si Action = delegate
	SubGoals  []string // si Action = split
	Message   string   // si Action = clarify o finish
}

func (s *Supervisor) Analyze(goal string, rt *AgentRuntime, ctx *AgentContext) SupervisorDecision {
	// 1. Goal vacío → aclarar
	if strings.TrimSpace(goal) == "" {
		return SupervisorDecision{
			Action:  "clarify",
			Message: "Necesito un objetivo para empezar.",
		}
	}

	// 2. Goal demasiado amplio → dividir
	if isTooBroad(goal) {
		return SupervisorDecision{
			Action:   "split",
			SubGoals: splitGoal(goal),
		}
	}

	// 3. Estancamiento detectado por el monitor del runtime
	if rt.Monitor != nil && rt.Monitor.ShouldReplan() {
		return SupervisorDecision{
			Action: "replan",
		}
	}

	// 4. Goal cumplido
	if isGoalSatisfied(goal, ctx) {
		return SupervisorDecision{
			Action:  "finish",
			Message: "Objetivo completado.",
		}
	}

	// 5. Usar memoria para las decisiones globales
	if rt.Memory.Recall("fail_count") != nil {
		if rt.Memory.Recall("fail_count").(int) >= 3 {
			return SupervisorDecision{
				Action:  "clarify",
				Message: "He detectado fallos repetidos. ¿Quieres redefinir el objetivo?",
			}
		}
	}

	meta := s.Metacognition.Evaluate(goal, rt, ctx)
	adj := s.Strategy.Adjust(meta, rt, ctx, goal)

	// Aplicar AOC
	s.AOC.Update(meta, rt, ctx)

	// Si la estrategia dice que hay que cambiar de agente
	if adj.ShouldSwitch {
		return SupervisorDecision{
			Action:    "delegate",
			AgentName: adj.SwitchTo,
		}
	}

	// Si la confianza es baja → pedir aclaración
	if meta.Confidence < 0.3 {
		return SupervisorDecision{
			Action:  "clarify",
			Message: "No estoy seguro de cómo proceder. " + meta.Advice,
		}
	}

	// Si hay loops o estancamiento → replanificar
	for _, f := range meta.Flags {
		if f == "loop_detected" || f == "stalled" {
			return SupervisorDecision{
				Action: "replan",
			}
		}
	}

	// Si hay fallos repetidos → dividir goal
	for _, f := range meta.Flags {
		if f == "repeated_failures" {
			return SupervisorDecision{
				Action:   "split",
				SubGoals: splitGoal(goal),
			}
		}
	}

	// 6. Delegación normal
	return SupervisorDecision{
		Action:    "delegate",
		AgentName: classifyGoal(goal, *s.llm),
	}
}

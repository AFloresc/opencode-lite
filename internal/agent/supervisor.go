package agent

import (
	"strings"
)

type Supervisor struct {
	llm *LLMClient
}

func NewSupervisor(llm LLMClient) *Supervisor {
	return &Supervisor{llm: &llm}
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

	// 5. Delegación normal
	return SupervisorDecision{
		Action:    "delegate",
		AgentName: classifyGoal(goal, *s.llm),
	}
}

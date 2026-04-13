package agent

type StrategyAdjustment struct {
	PlannerMode   string // "coarse", "fine", "aggressive", "conservative"
	GroundingMode string // "strict", "flexible"
	ToolBias      string // "analysis", "refactor", "docs", "none"
	ShouldSwitch  bool   // cambiar de agente
	SwitchTo      string // nombre del agente
}

type StrategyEngine struct {
	llm LLMClient
}

func NewStrategyEngine(llm LLMClient) *StrategyEngine {
	return &StrategyEngine{llm: llm}
}

func (s *StrategyEngine) Adjust(meta MetaEvaluation, rt *AgentRuntime, ctx *AgentContext, goal string) StrategyAdjustment {
	adj := StrategyAdjustment{
		PlannerMode:   "coarse",
		GroundingMode: "strict",
		ToolBias:      "none",
		ShouldSwitch:  false,
	}

	// ============================================================
	// 1. Baja confianza → estrategia conservadora
	// ============================================================
	if meta.Confidence < 0.3 {
		adj.PlannerMode = "fine"
		adj.GroundingMode = "strict"
		return adj
	}

	// ============================================================
	// 2. Loops → grounding flexible + planner agresivo
	// ============================================================
	for _, f := range meta.Flags {
		if f == "loop_detected" {
			adj.PlannerMode = "aggressive"
			adj.GroundingMode = "flexible"
			return adj
		}
	}

	// ============================================================
	// 3. Fallos repetidos → cambiar de agente
	// ============================================================
	for _, f := range meta.Flags {
		if f == "repeated_failures" {
			adj.ShouldSwitch = true
			adj.SwitchTo = classifyGoal(goal, s.llm)
			return adj
		}
	}

	// ============================================================
	// 4. Estancamiento → planner más granular
	// ============================================================
	for _, f := range meta.Flags {
		if f == "stalled" {
			adj.PlannerMode = "fine"
			return adj
		}
	}

	// ============================================================
	// 5. Patrones desde memoria → sesgo de tools
	// ============================================================
	if rt.Memory.Recall("dependency_cycles") == true {
		adj.ToolBias = "analysis"
	}

	if rt.Memory.Recall("long_functions_detected") == true {
		adj.ToolBias = "refactor"
	}

	return adj
}

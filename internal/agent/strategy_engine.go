package agent

//
// ============================================================
//  StrategyAdjustment
//  - indica cómo debe adaptarse el runtime cognitivo
// ============================================================
//

type StrategyAdjustment struct {
	PlannerMode   string // "coarse", "fine", "aggressive", "conservative"
	GroundingMode string // "strict", "flexible"
	ToolBias      string // "analysis", "refactor", "docs", "none"
	ShouldSwitch  bool   // cambiar de agente
	SwitchTo      string // nombre del agente
}

//
// ============================================================
//  StrategyEngine
//  - decide ajustes cognitivos basados en:
//      • metacognición
//      • memoria cognitiva
//      • señales del monitor
//      • tipo de goal
// ============================================================
//

type StrategyEngine struct {
	llm LLMClient
}

func NewStrategyEngine(llm LLMClient) *StrategyEngine {
	return &StrategyEngine{llm: llm}
}

//
// ============================================================
//  Adjust: núcleo de la estrategia adaptativa
// ============================================================
//

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
	// 2. Flags del monitor → loops, fallos, estancamiento
	// ============================================================

	for _, f := range meta.Flags {

		// Loop → grounding flexible + planner agresivo
		if f == "loop_detected" {
			adj.PlannerMode = "aggressive"
			adj.GroundingMode = "flexible"
			return adj
		}

		// Fallos repetidos → cambiar de agente
		if f == "repeated_failures" {
			adj.ShouldSwitch = true
			adj.SwitchTo = classifyGoal(goal, s.llm)
			return adj
		}

		// Estancamiento → granularidad fina
		if f == "stalled" {
			adj.PlannerMode = "fine"
			adj.GroundingMode = "flexible"
			return adj
		}
	}

	// ============================================================
	// 3. Sesgos basados en memoria cognitiva
	// ============================================================

	// Si se detectaron ciclos de dependencias → análisis
	if rt.Memory.Recall("dependency_cycles") == true {
		adj.ToolBias = "analysis"
	}

	// Si se detectaron funciones largas → refactor
	if rt.Memory.Recall("long_functions_detected") == true {
		adj.ToolBias = "refactor"
	}

	// Si hubo éxito reciente → modo coarse (más rápido)
	if success, ok := rt.Memory.Recall("success_count").(int); ok && success >= 3 {
		adj.PlannerMode = "coarse"
	}

	// Si hubo fallos recientes → modo fine
	if fails, ok := rt.Memory.Recall("fail_count").(int); ok && fails >= 2 {
		adj.PlannerMode = "fine"
	}

	// ============================================================
	// 4. Ajustes según tipo de goal
	// ============================================================

	g := classifyGoal(goal, s.llm)

	switch g {
	case "analysis":
		adj.PlannerMode = "fine"
		adj.ToolBias = "analysis"

	case "refactor":
		adj.PlannerMode = "coarse"
		adj.ToolBias = "refactor"

	case "docs":
		adj.PlannerMode = "fine"
		adj.ToolBias = "docs"
	}

	// ============================================================
	// 5. Ajustes según historial de herramientas
	// ============================================================

	if lastTool, ok := rt.Memory.Recall("last_tool").(string); ok {
		if lastTool == "search_in_project" {
			adj.GroundingMode = "flexible"
		}
	}

	// ============================================================
	// 6. Ajustes según tamaño del proyecto
	// ============================================================

	if rt.Grounder != nil {
		if stats, ok := rt.Grounder.(*ContextualToolGrounder); ok {
			if stats.Stats.FileCount > 500 {
				adj.PlannerMode = "coarse" // proyectos grandes → pasos más amplios
			}
		}
	}

	return adj
}

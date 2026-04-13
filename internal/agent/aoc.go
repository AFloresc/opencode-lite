package agent

import (
	"time"
)

//
// ============================================================
//  AOCUpdate
//  - snapshot persistente de preferencias aprendidas
// ============================================================
//

type AOCUpdate struct {
	PreferredPlannerMode   string    `json:"preferred_planner_mode"`
	PreferredGroundingMode string    `json:"preferred_grounding_mode"`
	PreferredAgentBias     string    `json:"preferred_agent_bias"`
	UpdatedAt              time.Time `json:"updated_at"`
}

//
// ============================================================
//  AOC (Auto‑Optimización Continua)
// ============================================================
//

type AOC struct {
	Memory *CognitiveMemory
}

func NewAOC(mem *CognitiveMemory) *AOC {
	return &AOC{Memory: mem}
}

//
// ============================================================
//  Update: núcleo del aprendizaje persistente
// ============================================================
//

func (a *AOC) Update(meta MetaEvaluation, rt *AgentRuntime, ctx *AgentContext) {
	update := AOCUpdate{
		PreferredPlannerMode:   a.computePlannerPreference(meta, rt),
		PreferredGroundingMode: a.computeGroundingPreference(meta, rt),
		PreferredAgentBias:     a.computeAgentBias(meta, rt),
		UpdatedAt:              time.Now(),
	}

	// Guardar snapshot
	a.Memory.Remember("aoc_last_update", update)

	// Guardar estadísticas acumuladas
	a.updateStats(meta, rt)

	a.Memory.Save()
}

//
// ============================================================
//  Preferencias aprendidas
// ============================================================
//

func (a *AOC) computePlannerPreference(meta MetaEvaluation, rt *AgentRuntime) string {

	// Baja confianza → granularidad fina
	if meta.Confidence < 0.3 {
		return "fine"
	}

	// Loops → agresivo
	if rt.Monitor.RepeatCount >= 2 {
		return "aggressive"
	}

	// Estancamiento → granular
	if rt.Monitor.StallCount >= 2 {
		return "fine"
	}

	// Fallos repetidos → conservador
	if rt.Monitor.FailureCount >= 2 {
		return "conservative"
	}

	// Éxito reciente → coarse
	if success, ok := rt.Memory.Recall("success_count").(int); ok && success >= 3 {
		return "coarse"
	}

	return "coarse"
}

func (a *AOC) computeGroundingPreference(meta MetaEvaluation, rt *AgentRuntime) string {

	// Loops → grounding flexible
	if rt.Monitor.RepeatCount >= 2 {
		return "flexible"
	}

	// Estancamiento → grounding flexible
	if rt.Monitor.StallCount >= 2 {
		return "flexible"
	}

	// Fallos repetidos → grounding estricto
	if rt.Monitor.FailureCount >= 2 {
		return "strict"
	}

	return "strict"
}

func (a *AOC) computeAgentBias(meta MetaEvaluation, rt *AgentRuntime) string {

	// Patrones detectados
	if rt.Memory.Recall("dependency_cycles") == true {
		return "analysis"
	}

	if rt.Memory.Recall("long_functions_detected") == true {
		return "refactor"
	}

	// Éxito reciente en documentación
	if lastTool, ok := rt.Memory.Recall("last_tool").(string); ok {
		if lastTool == "summarize_text" || lastTool == "explain_code" {
			return "docs"
		}
	}

	return "none"
}

//
// ============================================================
//  Estadísticas acumuladas para aprendizaje a largo plazo
// ============================================================
//

func (a *AOC) updateStats(meta MetaEvaluation, rt *AgentRuntime) {

	// Incrementar contadores cognitivos
	if rt.Monitor.FailureCount > 0 {
		a.Memory.Increment("aoc_failures_total")
	}

	if rt.Monitor.RepeatCount > 0 {
		a.Memory.Increment("aoc_loops_total")
	}

	if rt.Monitor.StallCount > 0 {
		a.Memory.Increment("aoc_stalls_total")
	}

	// Guardar última confianza
	a.Memory.Remember("aoc_last_confidence", meta.Confidence)
}

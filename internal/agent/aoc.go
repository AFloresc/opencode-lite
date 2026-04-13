package agent

import (
	"time"
)

type AOC struct {
	Memory *CognitiveMemory
}

func NewAOC(mem *CognitiveMemory) *AOC {
	return &AOC{Memory: mem}
}

type AOCUpdate struct {
	PreferredPlannerMode   string
	PreferredGroundingMode string
	PreferredAgentBias     string
	UpdatedAt              time.Time
}

func (a *AOC) Update(meta MetaEvaluation, rt *AgentRuntime, ctx *AgentContext) {
	update := AOCUpdate{
		PreferredPlannerMode:   a.computePlannerPreference(meta, rt),
		PreferredGroundingMode: a.computeGroundingPreference(meta, rt),
		PreferredAgentBias:     a.computeAgentBias(meta, rt),
		UpdatedAt:              time.Now(),
	}

	a.Memory.Remember("aoc_update", update)
	a.Memory.Save()
}

func (a *AOC) computePlannerPreference(meta MetaEvaluation, rt *AgentRuntime) string {
	if meta.Confidence < 0.3 {
		return "fine"
	}
	if rt.Monitor.RepeatCount >= 2 {
		return "aggressive"
	}
	if rt.Monitor.StallCount >= 2 {
		return "fine"
	}
	return "coarse"
}

func (a *AOC) computeGroundingPreference(meta MetaEvaluation, rt *AgentRuntime) string {
	if rt.Monitor.RepeatCount >= 2 {
		return "flexible"
	}
	return "strict"
}

func (a *AOC) computeAgentBias(meta MetaEvaluation, rt *AgentRuntime) string {
	if rt.Memory.Recall("dependency_cycles") == true {
		return "analysis"
	}
	if rt.Memory.Recall("long_functions_detected") == true {
		return "refactor"
	}
	return "none"
}

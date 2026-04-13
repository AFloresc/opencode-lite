package agent

import (
	"encoding/json"
	"os"
	"strings"
)

type PlannerMemory struct {
	SuccessfulSteps map[string]int `json:"successful_steps"`
	FailedSteps     map[string]int `json:"failed_steps"`
	GoalPatterns    map[string]int `json:"goal_patterns"`
}

func NewPlannerMemory() *PlannerMemory {
	return &PlannerMemory{
		SuccessfulSteps: map[string]int{},
		FailedSteps:     map[string]int{},
		GoalPatterns:    map[string]int{},
	}
}

//
// ============================
// PERSISTENCIA EN DISCO
// ============================
//

const memoryFile = "agent_memory.json"

func (m *PlannerMemory) Save() error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(memoryFile, data, 0644)
}

func (m *PlannerMemory) Load() error {
	data, err := os.ReadFile(memoryFile)
	if err != nil {
		// Si no existe, no es error
		return nil
	}
	return json.Unmarshal(data, m)
}

//
// ============================
// ACTUALIZACIÓN DE MEMORIA
// ============================
//

func (m *PlannerMemory) RecordSuccess(step string) {
	m.SuccessfulSteps[strings.ToLower(step)]++
}

func (m *PlannerMemory) RecordFailure(step string) {
	m.FailedSteps[strings.ToLower(step)]++
}

func (m *PlannerMemory) RecordGoal(goal string) {
	m.GoalPatterns[strings.ToLower(goal)]++
}

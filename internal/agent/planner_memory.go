package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

//
// ============================================================
//  PlannerMemory
//  - memoria persistente del planner
//  - registra:
//      • pasos exitosos
//      • pasos fallidos
//      • patrones de goals
//  - usada por HybridPlanner para ordenar y filtrar pasos
// ============================================================
//

type PlannerMemory struct {
	SuccessfulSteps map[string]int `json:"successful_steps"`
	FailedSteps     map[string]int `json:"failed_steps"`
	GoalPatterns    map[string]int `json:"goal_patterns"`
	ProjectID       string         `json:"project_id"`
}

func NewPlannerMemory(projectID string) *PlannerMemory {
	return &PlannerMemory{
		SuccessfulSteps: map[string]int{},
		FailedSteps:     map[string]int{},
		GoalPatterns:    map[string]int{},
		ProjectID:       projectID,
	}
}

//
// ============================================================
//  Path
// ============================================================
//

func (m *PlannerMemory) memoryFilePath() string {
	dir := filepath.Join(".opencode", m.ProjectID)
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "planner_memory.json")
}

//
// ============================================================
//  Load / Save
// ============================================================
//

func (m *PlannerMemory) Load() error {
	data, err := os.ReadFile(m.memoryFilePath())
	if err != nil {
		return nil // no existe → memoria vacía
	}
	return json.Unmarshal(data, m)
}

func (m *PlannerMemory) Save() error {
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.memoryFilePath(), bytes, 0644)
}

//
// ============================================================
//  Recorders
// ============================================================
//

func (m *PlannerMemory) RecordSuccess(step string) {
	step = strings.ToLower(strings.TrimSpace(step))
	if step == "" {
		return
	}
	m.SuccessfulSteps[step]++
}

func (m *PlannerMemory) RecordFailure(step string) {
	step = strings.ToLower(strings.TrimSpace(step))
	if step == "" {
		return
	}
	m.FailedSteps[step]++
}

func (m *PlannerMemory) RecordGoal(goal string) {
	goal = strings.ToLower(strings.TrimSpace(goal))
	if goal == "" {
		return
	}
	m.GoalPatterns[goal]++
}

//
// ============================================================
//  Cognitive Helpers
// ============================================================
//

// Devuelve true si un step ha fallado demasiadas veces
func (m *PlannerMemory) IsFrequentlyFailing(step string) bool {
	step = strings.ToLower(strings.TrimSpace(step))
	return m.FailedSteps[step] >= 3
}

// Score cognitivo para ordenar pasos
func (m *PlannerMemory) Score(step string) int {
	step = strings.ToLower(strings.TrimSpace(step))
	return m.SuccessfulSteps[step] - m.FailedSteps[step]
}

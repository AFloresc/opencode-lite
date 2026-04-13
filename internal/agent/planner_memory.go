package agent

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

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

func (m *PlannerMemory) memoryFilePath() string {
	// Carpeta local oculta para memorias
	dir := ".agent_memory"

	_ = os.MkdirAll(dir, 0755)

	// Hash del projectID para evitar nombres raros
	h := sha1.Sum([]byte(m.ProjectID))
	filename := hex.EncodeToString(h[:]) + ".json"

	return filepath.Join(dir, filename)
}

func (m *PlannerMemory) Save() error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.memoryFilePath(), data, 0644)
}

func (m *PlannerMemory) Load() error {
	data, err := os.ReadFile(m.memoryFilePath())
	if err != nil {
		// Si no existe, no es error
		return nil
	}
	return json.Unmarshal(data, m)
}

func (m *PlannerMemory) RecordSuccess(step string) {
	m.SuccessfulSteps[strings.ToLower(step)]++
}

func (m *PlannerMemory) RecordFailure(step string) {
	m.FailedSteps[strings.ToLower(step)]++
}

func (m *PlannerMemory) RecordGoal(goal string) {
	m.GoalPatterns[strings.ToLower(goal)]++
}

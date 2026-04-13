package agent

import "strings"

type PlannerMemory struct {
	SuccessfulSteps map[string]int
	FailedSteps     map[string]int
	GoalPatterns    map[string]int
}

func NewPlannerMemory() *PlannerMemory {
	return &PlannerMemory{
		SuccessfulSteps: map[string]int{},
		FailedSteps:     map[string]int{},
		GoalPatterns:    map[string]int{},
	}
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

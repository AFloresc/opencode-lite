package agent

import (
	"opencode-lite/internal/tools"
	"time"
)

type ExecutionMonitor struct {
	LastTool       string
	RepeatCount    int
	FailureCount   int
	StallCount     int
	LastResultHash string
	LastStepTime   time.Time
}

func NewExecutionMonitor() *ExecutionMonitor {
	return &ExecutionMonitor{
		LastStepTime: time.Now(),
	}
}

func (m *ExecutionMonitor) Update(step PlanStep, result tools.ToolResult) {
	// 1. Detectar repetición de la misma tool
	if step.Description == m.LastTool {
		m.RepeatCount++
	} else {
		m.RepeatCount = 0
	}

	m.LastTool = step.Description

	// 2. Detectar fallos repetidos
	if result.Error != "" {
		m.FailureCount++
	} else {
		m.FailureCount = 0
	}

	// 3. Detectar estancamiento (mismo resultado)
	hash := hashResult(result)
	if hash == m.LastResultHash {
		m.StallCount++
	} else {
		m.StallCount = 0
	}
	m.LastResultHash = hash

	m.LastStepTime = time.Now()
}

func (m *ExecutionMonitor) ShouldReplan() bool {
	// 1. Loop evidente
	if m.RepeatCount >= 3 {
		return true
	}

	// 2. Fallos repetidos
	if m.FailureCount >= 2 {
		return true
	}

	// 3. Estancamiento
	if m.StallCount >= 3 {
		return true
	}

	return false
}

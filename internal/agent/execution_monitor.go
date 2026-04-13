package agent

import (
	"opencode-lite/internal/tools"
	"strings"
	"time"
)

type ExecutionMonitor struct {
	LastStepDesc   string
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

//
// ============================================================
//  Update: se llama después de cada tool execution
// ============================================================
//

func (m *ExecutionMonitor) Update(step PlanStep, result tools.ToolResult) {
	desc := strings.ToLower(step.Description)

	// 1. Detectar repetición del mismo step (loop)
	if desc == m.LastStepDesc {
		m.RepeatCount++
	} else {
		m.RepeatCount = 0
	}
	m.LastStepDesc = desc

	// 2. Detectar fallos repetidos
	if result.Error != "" {
		m.FailureCount++
	} else {
		m.FailureCount = 0
	}

	// 3. Detectar estancamiento por resultado idéntico
	hash := hashResult(result)
	if hash == m.LastResultHash {
		m.StallCount++
	} else {
		m.StallCount = 0
	}
	m.LastResultHash = hash

	// 4. Actualizar timestamp
	m.LastStepTime = time.Now()
}

//
// ============================================================
//  ShouldReplan: decide si hay que replanificar
// ============================================================
//

func (m *ExecutionMonitor) ShouldReplan() bool {

	// 1. Loop evidente
	if m.RepeatCount >= 3 {
		return true
	}

	// 2. Fallos repetidos
	if m.FailureCount >= 2 {
		return true
	}

	// 3. Estancamiento por resultado idéntico
	if m.StallCount >= 3 {
		return true
	}

	// 4. Estancamiento temporal (sin progreso)
	if time.Since(m.LastStepTime) > 4*time.Second {
		return true
	}

	return false
}

//
// ============================================================
//  Flags: señales para Metacognición y Supervisor
// ============================================================
//

func (m *ExecutionMonitor) Flags() []string {
	flags := []string{}

	if m.RepeatCount >= 3 {
		flags = append(flags, "loop_detected")
	}

	if m.FailureCount >= 2 {
		flags = append(flags, "repeated_failures")
	}

	if m.StallCount >= 3 || time.Since(m.LastStepTime) > 4*time.Second {
		flags = append(flags, "stalled")
	}

	return flags
}

package agent

import (
	"sort"
	"strings"
)

type PlanStep struct {
	Description string
}

type Plan struct {
	Steps []PlanStep
}

type Planner interface {
	MakePlan(goal string) Plan
	UpdateMemory(ctx AgentContext)
}

type MemoryPlanner struct {
	Memory *PlannerMemory
}

func NewMemoryPlanner() *MemoryPlanner {
	return &MemoryPlanner{
		Memory: NewPlannerMemory(),
	}
}

func (p *MemoryPlanner) MakePlan(goal string) Plan {
	g := strings.ToLower(goal)
	p.Memory.RecordGoal(g)

	// Base plan (como el planner avanzado anterior)
	plan := p.basePlan(g)

	// Ajustar plan según memoria
	plan = p.applyMemoryHeuristics(plan)

	return plan
}

//
// ============================
// PLAN BASE (igual que el avanzado anterior)
// ============================
//

func (p *MemoryPlanner) basePlan(goal string) Plan {
	switch {
	case containsAny(goal, "analizar proyecto", "project analysis"):
		return Plan{
			Steps: []PlanStep{
				{"listar archivos"},
				{"calcular métricas"},
				{"detectar dependencias"},
				{"buscar duplicación"},
				{"buscar funciones largas"},
			},
		}

	case containsAny(goal, "explicar archivo", "explain file"):
		return Plan{
			Steps: []PlanStep{
				{"extraer funciones"},
				{"extraer tipos"},
				{"extraer comentarios"},
				{"resumir archivo"},
			},
		}

	case containsAny(goal, "limpiar proyecto", "clean project"):
		return Plan{
			Steps: []PlanStep{
				{"listar archivos"},
				{"limpiar imports"},
				{"formatear"},
			},
		}
	}

	return Plan{Steps: []PlanStep{{goal}}}
}

//
// ============================
// HEURÍSTICAS BASADAS EN MEMORIA
// ============================
//

func (p *MemoryPlanner) applyMemoryHeuristics(plan Plan) Plan {
	// Ordenar pasos según éxito histórico
	sort.Slice(plan.Steps, func(i, j int) bool {
		a := strings.ToLower(plan.Steps[i].Description)
		b := strings.ToLower(plan.Steps[j].Description)

		return p.Memory.SuccessfulSteps[a] > p.Memory.SuccessfulSteps[b]
	})

	// Eliminar pasos que fallaron muchas veces
	filtered := []PlanStep{}
	for _, step := range plan.Steps {
		if p.Memory.FailedSteps[strings.ToLower(step.Description)] < 3 {
			filtered = append(filtered, step)
		}
	}

	plan.Steps = filtered
	return plan
}

//
// ============================
// ACTUALIZAR MEMORIA DESPUÉS DE EJECUTAR
// ============================
//

func (p *MemoryPlanner) UpdateMemory(ctx AgentContext) {
	for _, step := range ctx.History {
		if step.Output.Error == "" {
			p.Memory.RecordSuccess(step.Action)
		} else {
			p.Memory.RecordFailure(step.Action)
		}
	}
}

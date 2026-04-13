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
	mem := NewPlannerMemory()
	mem.Load() // ← carga memoria persistente

	return &MemoryPlanner{
		Memory: mem,
	}
}

func (p *MemoryPlanner) MakePlan(goal string) Plan {
	g := strings.ToLower(goal)
	p.Memory.RecordGoal(g)

	plan := p.basePlan(g)
	plan = p.applyMemoryHeuristics(plan)

	return plan
}

//
// ============================
// PLAN BASE
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
	}

	return Plan{Steps: []PlanStep{{goal}}}
}

//
// ============================
// HEURÍSTICAS BASADAS EN MEMORIA
// ============================
//

func (p *MemoryPlanner) applyMemoryHeuristics(plan Plan) Plan {
	// Ordenar por éxito histórico
	sort.Slice(plan.Steps, func(i, j int) bool {
		a := strings.ToLower(plan.Steps[i].Description)
		b := strings.ToLower(plan.Steps[j].Description)
		return p.Memory.SuccessfulSteps[a] > p.Memory.SuccessfulSteps[b]
	})

	// Eliminar pasos que fallaron repetidamente
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
// ACTUALIZAR MEMORIA Y GUARDAR
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

	p.Memory.Save() // ← persistencia automática
}

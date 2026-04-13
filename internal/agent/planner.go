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

type HybridPlanner struct {
	Memory  *PlannerMemory
	LLM     LLMClient
	Project string
	Mode    string
}

func NewHybridPlanner(projectID string, llm LLMClient) *HybridPlanner {
	mem := NewPlannerMemory(projectID)
	_ = mem.Load()

	return &HybridPlanner{
		Memory:  mem,
		LLM:     llm,
		Project: projectID,
		Mode:    "coarse", // default mode
	}
}

func (p *HybridPlanner) MakePlan(goal string) Plan {
	g := strings.ToLower(goal)
	p.Memory.RecordGoal(g)

	// 1) Intentar reglas rápidas
	if plan, ok := matchQuickRule(g); ok {
		plan = p.applyMemoryHeuristics(plan)
		return plan
	}

	// 2) Intentar reglas base (más generales)
	plan := p.basePlan(g)

	// 3) Si el plan es trivial → LLM
	if len(plan.Steps) == 1 && strings.EqualFold(plan.Steps[0].Description, goal) && p.LLM != nil {
		if llmPlan, ok := p.tryLLMPlan(goal); ok {
			plan = llmPlan
		}
	}

	// 4) Aplicar heurísticas de memoria
	plan = p.applyMemoryHeuristics(plan)

	return plan
}

func (p *HybridPlanner) tryLLMPlan(goal string) (Plan, bool) {
	steps, err := p.LLM.ProposePlan(goal)
	if err != nil || len(steps) == 0 {
		return Plan{}, false
	}

	plan := Plan{Steps: make([]PlanStep, 0, len(steps))}
	for _, s := range steps {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		plan.Steps = append(plan.Steps, PlanStep{Description: s})
	}

	if len(plan.Steps) == 0 {
		return Plan{}, false
	}

	return plan, true
}

//
// ============================
// PLAN BASE (REGLAS)
// ============================
//

func (p *HybridPlanner) basePlan(goal string) Plan {
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

	// Plan trivial → candidato a LLM
	return Plan{Steps: []PlanStep{{Description: goal}}}
}

//
// ============================
// HEURÍSTICAS BASADAS EN MEMORIA
// ============================
//

func (p *HybridPlanner) applyMemoryHeuristics(plan Plan) Plan {
	sort.Slice(plan.Steps, func(i, j int) bool {
		a := strings.ToLower(plan.Steps[i].Description)
		b := strings.ToLower(plan.Steps[j].Description)
		return p.Memory.SuccessfulSteps[a] > p.Memory.SuccessfulSteps[b]
	})

	filtered := []PlanStep{}
	for _, step := range plan.Steps {
		if p.Memory.FailedSteps[strings.ToLower(step.Description)] < 3 {
			filtered = append(filtered, step)
		}
	}

	plan.Steps = filtered
	return plan
}

func (p *HybridPlanner) UpdateMemory(ctx AgentContext) {
	for _, step := range ctx.History {
		if step.Output.Error == "" {
			p.Memory.RecordSuccess(step.Action)
		} else {
			p.Memory.RecordFailure(step.Action)
		}
	}

	_ = p.Memory.Save()
}

func (p *HybridPlanner) SetMode(mode string) {
	p.Mode = mode
}

func (p *HybridPlanner) applyPlannerMode(plan Plan) Plan {
	switch p.Mode {

	case "fine":
		return p.expandFine(plan)

	case "aggressive":
		return p.expandAggressive(plan)

	case "conservative":
		return p.simplifyConservative(plan)

	case "coarse":
		fallthrough
	default:
		return plan
	}
}

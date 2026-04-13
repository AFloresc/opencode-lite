package agent

func (p *HybridPlanner) expandFine(plan Plan) Plan {
	out := Plan{}
	for _, step := range plan.Steps {
		if len(step.Description) > 40 {
			out.Steps = append(out.Steps,
				PlanStep{Description: "analizar: " + step.Description},
				PlanStep{Description: "ejecutar: " + step.Description},
			)
		} else {
			out.Steps = append(out.Steps, step)
		}
	}
	return out
}

func (p *HybridPlanner) expandAggressive(plan Plan) Plan {
	if len(plan.Steps) == 0 {
		return plan
	}
	out := plan
	out.Steps = append([]PlanStep{
		{Description: "explorar alternativas para: " + plan.Steps[0].Description},
	}, out.Steps...)
	return out
}

func (p *HybridPlanner) simplifyConservative(plan Plan) Plan {
	if len(plan.Steps) > 1 {
		return Plan{Steps: []PlanStep{plan.Steps[0]}}
	}
	return plan
}

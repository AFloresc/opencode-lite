package agent

type AgentRuntime struct {
	Policy   AgentPolicy
	Planner  Planner
	Mapper   StepMapper
	Grounder ToolGrounder
	Expander StepExpander
}

func NewAgentRuntime(projectID string, policy AgentPolicy, llm LLMClient) *AgentRuntime {
	stats := AnalyzeProjectSize()
	mem := NewPlannerMemory(projectID)
	_ = mem.Load()

	return &AgentRuntime{
		Policy:   policy,
		Planner:  NewHybridPlanner(projectID, llm),
		Mapper:   NewSemanticStepMapper(),
		Grounder: NewContextualToolGrounder(stats, mem),
		Expander: NewDefaultStepExpander(),
	}
}

func (rt *AgentRuntime) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	plan := rt.Planner.MakePlan(goal)
	queue := append([]PlanStep{}, plan.Steps...)

	for len(queue) > 0 {
		step := queue[0]
		queue = queue[1:]

		normalized := rt.Mapper.Normalize(step.Description)
		ctx.Goal = normalized

		// 1. Grounding directo
		if call, ok := rt.Grounder.Ground(normalized, &ctx); ok {
			result := rt.executeTool(call.ToolName, call.Args, &ctx)

			// 2. Subplanes dinámicos
			if rt.Expander != nil {
				newSteps := rt.Expander.Expand(step, result, &ctx)
				queue = append(newSteps, queue...)
			}

			continue
		}

		// 3. Fallback: Policy
		for i := 0; i < 20; i++ {
			toolName, args, done := rt.Policy.Decide(&ctx)
			if done {
				break
			}

			result := rt.executeTool(toolName, args, &ctx)

			if rt.Expander != nil {
				newSteps := rt.Expander.Expand(step, result, &ctx)
				queue = append(newSteps, queue...)
			}
		}
	}

	rt.Planner.UpdateMemory(ctx)
	return ctx
}

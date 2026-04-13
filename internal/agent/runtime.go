package agent

type AgentRuntime struct {
	Policy   AgentPolicy
	Planner  Planner
	Mapper   StepMapper
	Grounder ToolGrounder
	Expander StepExpander
	Monitor  *ExecutionMonitor
	Memory   *CognitiveMemory
}

func NewAgentRuntime(projectID string, policy AgentPolicy, llm LLMClient) *AgentRuntime {
	stats := AnalyzeProjectSize()
	memPlanner := NewPlannerMemory(projectID)
	_ = memPlanner.Load()
	mem := NewCognitiveMemory(projectID)
	_ = mem.Load()

	return &AgentRuntime{
		Policy:   policy,
		Planner:  NewHybridPlanner(projectID, llm),
		Mapper:   NewSemanticStepMapper(),
		Grounder: NewContextualToolGrounder(stats, mem),
		Expander: NewDefaultStepExpander(),
		Monitor:  NewExecutionMonitor(),
		Memory:   mem,
	}
}

func (rt *AgentRuntime) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	// 1. Generar plan inicial
	plan := rt.Planner.MakePlan(goal)
	queue := append([]PlanStep{}, plan.Steps...)

	for len(queue) > 0 {
		step := queue[0]
		queue = queue[1:]

		// Normalizamos el step
		normalized := rt.Mapper.Normalize(step.Description)
		ctx.Goal = normalized

		// ============================================================
		// 2. Grounding contextual (primer intento)
		// ============================================================
		if call, ok := rt.Grounder.Ground(normalized, &ctx); ok {

			result := rt.executeTool(call.ToolName, call.Args, &ctx)

			// Guardar en memoria cognitiva
			rt.Memory.Remember("last_tool", call.ToolName)
			rt.Memory.Remember("last_result", result.Result)
			rt.Memory.Save()

			// Monitor de ejecución
			rt.Monitor.Update(step, result)

			// Subplanes dinámicos
			if rt.Expander != nil {
				newSteps := rt.Expander.Expand(step, result, &ctx)
				queue = append(newSteps, queue...)
			}

			// Replanificación automática
			if rt.Monitor.ShouldReplan() {
				newPlan := rt.Planner.MakePlan(ctx.Goal)
				queue = append(newPlan.Steps, queue...)
				rt.Monitor = NewExecutionMonitor()
			}

			continue
		}

		// ============================================================
		// 3. Fallback: Policy Decide()
		// ============================================================
		for i := 0; i < 20; i++ {
			toolName, args, done := rt.Policy.Decide(&ctx)
			if done {
				break
			}

			result := rt.executeTool(toolName, args, &ctx)

			// Guardar en memoria cognitiva
			rt.Memory.Remember("last_tool", toolName)
			rt.Memory.Remember("last_result", result.Result)
			rt.Memory.Save()

			// Monitor
			rt.Monitor.Update(step, result)

			// Subplanes dinámicos
			if rt.Expander != nil {
				newSteps := rt.Expander.Expand(step, result, &ctx)
				queue = append(newSteps, queue...)
			}

			// Replanificación automática
			if rt.Monitor.ShouldReplan() {
				newPlan := rt.Planner.MakePlan(ctx.Goal)
				queue = append(newPlan.Steps, queue...)
				rt.Monitor = NewExecutionMonitor()
			}
		}
	}

	// Guardar memoria del planner
	rt.Planner.UpdateMemory(ctx)

	return ctx
}

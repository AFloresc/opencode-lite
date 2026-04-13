package agent

import (
	"opencode-lite/internal/tools"
)

type AgentRuntime struct {
	Policy   AgentPolicy
	Planner  Planner
	Mapper   StepMapper
	Grounder ToolGrounder
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
	}
}

func (rt *AgentRuntime) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	plan := rt.Planner.MakePlan(goal)

	for _, step := range plan.Steps {
		normalized := step.Description
		if rt.Mapper != nil {
			normalized = rt.Mapper.Normalize(step.Description)
		}

		ctx.Goal = normalized

		// 1) Intentar grounding directo
		if rt.Grounder != nil {
			if call, ok := rt.Grounder.Ground(normalized, &ctx); ok {
				toolFn, ok := tools.ToolRegistry[call.ToolName]
				if !ok {
					ctx.LastResult = tools.ToolResult{
						ToolName: call.ToolName,
						Error:    "tool no encontrada",
					}
					continue
				}

				result := toolFn(call.Args)

				ctx.History = append(ctx.History, AgentStep{
					Thought: "Ejecutando " + call.ToolName + " (grounded)",
					Action:  call.ToolName,
					Input:   call.Args,
					Output:  result,
				})

				ctx.LastResult = result
				continue
			}
		}

		// 2) Fallback: usar Policy
		for i := 0; i < 20; i++ {
			toolName, args, done := rt.Policy.Decide(&ctx)
			if done {
				break
			}

			toolFn, ok := tools.ToolRegistry[toolName]
			if !ok {
				ctx.LastResult = tools.ToolResult{
					ToolName: toolName,
					Error:    "tool no encontrada",
				}
				break
			}

			result := toolFn(args)

			ctx.History = append(ctx.History, AgentStep{
				Thought: "Ejecutando " + toolName,
				Action:  toolName,
				Input:   args,
				Output:  result,
			})

			ctx.LastResult = result
		}
	}

	rt.Planner.UpdateMemory(ctx)
	return ctx
}

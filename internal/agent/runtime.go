package agent

import (
	"opencode-lite/internal/tools"
)

type AgentRuntime struct {
	Policy  AgentPolicy
	Planner Planner
}

func NewAgentRuntime(projectID string, policy AgentPolicy) *AgentRuntime {
	return &AgentRuntime{
		Policy:  policy,
		Planner: NewMemoryPlanner(projectID),
	}
}

func (rt *AgentRuntime) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	plan := rt.Planner.MakePlan(goal)

	for _, step := range plan.Steps {
		ctx.Goal = step.Description

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

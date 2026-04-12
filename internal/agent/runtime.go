package agent

import "opencode-lite/internal/tools"

type AgentRuntime struct {
	Policy AgentPolicy
}

func (rt *AgentRuntime) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	for step := 0; step < 20; step++ {
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

	return ctx
}

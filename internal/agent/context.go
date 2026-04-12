package agent

import "opencode-lite/internal/tools"

type AgentContext struct {
	Goal       string
	Memory     map[string]interface{}
	LastResult tools.ToolResult
	History    []AgentStep
}

type AgentStep struct {
	Thought string
	Action  string
	Input   map[string]interface{}
	Output  tools.ToolResult
}

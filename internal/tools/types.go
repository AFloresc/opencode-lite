package tools

// ToolCall representa una llamada a herramienta generada por el modelo
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult representa el resultado de ejecutar una herramienta
type ToolResult struct {
	ToolName string      `json:"tool_name"`
	Result   interface{} `json:"result"`
	Error    string      `json:"error,omitempty"`
}

package agent

type GroundedToolCall struct {
	ToolName string
	Args     map[string]interface{}
}

type ToolGrounder interface {
	Ground(step string, ctx *AgentContext) (*GroundedToolCall, bool)
}

type DefaultToolGrounder struct{}

func NewDefaultToolGrounder() *DefaultToolGrounder {
	return &DefaultToolGrounder{}
}

func (g *DefaultToolGrounder) Ground(step string, ctx *AgentContext) (*GroundedToolCall, bool) {
	switch step {
	case "listar archivos":
		return &GroundedToolCall{
			ToolName: "file_tree",
			Args: map[string]interface{}{
				"root": "workspace",
			},
		}, true

	case "calcular métricas":
		return &GroundedToolCall{
			ToolName: "analysis_metrics",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "detectar dependencias":
		return &GroundedToolCall{
			ToolName: "analysis_dependencies",
			Args: map[string]interface{}{
				"root": "workspace",
			},
		}, true

	case "buscar duplicación":
		return &GroundedToolCall{
			ToolName: "search_regex_multi",
			Args: map[string]interface{}{
				"path":    "workspace",
				"pattern": "func .*\\{[\\s\\S]{100,}\\}",
			},
		}, true

	case "buscar funciones largas":
		return &GroundedToolCall{
			ToolName: "search_regex_multi",
			Args: map[string]interface{}{
				"path":    "workspace",
				"pattern": "func .*\\{[\\s\\S]{200,}\\}",
			},
		}, true

	case "limpiar imports":
		return &GroundedToolCall{
			ToolName: "format_code",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "formatear":
		return &GroundedToolCall{
			ToolName: "format_code",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "extraer funciones":
		return &GroundedToolCall{
			ToolName: "extract_functions",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "extraer tipos":
		return &GroundedToolCall{
			ToolName: "extract_types",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "extraer comentarios":
		return &GroundedToolCall{
			ToolName: "extract_comments_block",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "resumir archivo":
		return &GroundedToolCall{
			ToolName: "summarize_file",
			Args: map[string]interface{}{
				"path": "workspace",
			},
		}, true

	case "dead code":
		return &GroundedToolCall{
			ToolName: "analysis_dead_code",
			Args: map[string]interface{}{
				"root": "workspace",
			},
		}, true
	}

	return nil, false
}

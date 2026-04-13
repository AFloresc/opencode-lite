package agent

import "strings"

type ContextualToolGrounder struct {
	ProjectStats ProjectStats
	Memory       *PlannerMemory
}

func NewContextualToolGrounder(stats ProjectStats, mem *PlannerMemory) *ContextualToolGrounder {
	return &ContextualToolGrounder{
		ProjectStats: stats,
		Memory:       mem,
	}
}

func (g *ContextualToolGrounder) Ground(step string, ctx *AgentContext) (*GroundedToolCall, bool) {
	s := strings.ToLower(step)

	// === 1. Si hay archivo actual, usarlo ===
	currentFile := ""
	if v, ok := ctx.Memory["current_file"].(string); ok {
		currentFile = v
	}

	// === 2. Si el proyecto es grande, usar herramientas rápidas ===
	isLarge := g.ProjectStats.FileCount > 300

	// === 3. Si la memoria dice que un tool falla mucho, evitarlo ===
	avoid := func(tool string) bool {
		return g.Memory.FailedSteps[strings.ToLower(tool)] > 5
	}

	// === 4. Grounding contextual ===

	switch s {

	case "listar archivos":
		if isLarge {
			return &GroundedToolCall{
				ToolName: "file_tree_fast",
				Args:     map[string]interface{}{"root": "workspace"},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "file_tree",
			Args:     map[string]interface{}{"root": "workspace"},
		}, true

	case "calcular métricas":
		if avoid("analysis_metrics") {
			return &GroundedToolCall{
				ToolName: "analysis_metrics_fast",
				Args:     map[string]interface{}{"path": "workspace"},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "analysis_metrics",
			Args:     map[string]interface{}{"path": "workspace"},
		}, true

	case "detectar dependencias":
		if avoid("analysis_dependencies") {
			return &GroundedToolCall{
				ToolName: "analysis_dependencies_light",
				Args:     map[string]interface{}{"root": "workspace"},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "analysis_dependencies",
			Args:     map[string]interface{}{"root": "workspace"},
		}, true

	case "extraer funciones":
		if currentFile != "" {
			return &GroundedToolCall{
				ToolName: "extract_functions",
				Args:     map[string]interface{}{"path": currentFile},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "extract_functions",
			Args:     map[string]interface{}{"path": "workspace"},
		}, true

	case "extraer tipos":
		if currentFile != "" {
			return &GroundedToolCall{
				ToolName: "extract_types",
				Args:     map[string]interface{}{"path": currentFile},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "extract_types",
			Args:     map[string]interface{}{"path": "workspace"},
		}, true

	case "extraer comentarios":
		if currentFile != "" {
			return &GroundedToolCall{
				ToolName: "extract_comments_block",
				Args:     map[string]interface{}{"path": currentFile},
			}, true
		}
		return &GroundedToolCall{
			ToolName: "extract_comments_block",
			Args:     map[string]interface{}{"path": "workspace"},
		}, true

	case "resumir archivo":
		if currentFile != "" {
			return &GroundedToolCall{
				ToolName: "summarize_file",
				Args:     map[string]interface{}{"path": currentFile},
			}, true
		}
	}
	return nil, false
}

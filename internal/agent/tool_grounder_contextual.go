package agent

import (
	"strings"
)

//
// ============================================================
//  ContextualToolGrounder
// ============================================================
//

type ContextualToolGrounder struct {
	Stats  ProjectStats
	Memory *CognitiveMemory
	Mode   string // "strict" o "flexible"
}

func NewContextualToolGrounder(stats ProjectStats, mem *CognitiveMemory) *ContextualToolGrounder {
	return &ContextualToolGrounder{
		Stats:  stats,
		Memory: mem,
		Mode:   "strict",
	}
}

func (g *ContextualToolGrounder) SetMode(mode string) {
	if mode != "" {
		g.Mode = mode
	}
}

//
// ============================================================
//  Grounding principal
// ============================================================
//

func (g *ContextualToolGrounder) Ground(step string, ctx *AgentContext) (*GroundedToolCall, bool) {
	s := strings.ToLower(step)

	// 1) Grounding estricto
	if g.Mode == "strict" {
		if call, ok := g.strictGround(s, ctx); ok {
			return call, true
		}
	}

	// 2) Grounding flexible
	if call, ok := g.flexibleGround(s, ctx); ok {
		return call, true
	}

	return nil, false
}

//
// ============================================================
//  STRICT GROUNDING
// ============================================================
//

func (g *ContextualToolGrounder) strictGround(step string, ctx *AgentContext) (*GroundedToolCall, bool) {

	// Dependencias
	if containsAny(step, "dependencia", "dependencias", "imports") {
		return &GroundedToolCall{
			ToolName: "list_dependencies",
			Args:     map[string]interface{}{},
		}, true
	}

	// Métricas
	if containsAny(step, "métrica", "métricas", "metrics") {
		return &GroundedToolCall{
			ToolName: "compute_metrics",
			Args:     map[string]interface{}{},
		}, true
	}

	// Listar archivos
	if containsAny(step, "listar archivos", "list files") {
		return &GroundedToolCall{
			ToolName: "list_files",
			Args:     map[string]interface{}{},
		}, true
	}

	// Formatear
	if containsAny(step, "formatear", "format") {
		return &GroundedToolCall{
			ToolName: "format_code",
			Args:     map[string]interface{}{},
		}, true
	}

	// Limpiar imports
	if containsAny(step, "limpiar imports", "clean imports") {
		return &GroundedToolCall{
			ToolName: "clean_imports",
			Args:     map[string]interface{}{},
		}, true
	}

	return nil, false
}

//
// ============================================================
//  FLEXIBLE GROUNDING
// ============================================================
//

func (g *ContextualToolGrounder) flexibleGround(step string, ctx *AgentContext) (*GroundedToolCall, bool) {

	// Si menciona archivo
	if containsAny(step, "archivo", "file") {
		file := g.extractFileName(step)
		if file != "" {
			return &GroundedToolCall{
				ToolName: "read_file",
				Args: map[string]interface{}{
					"path": file,
				},
			}, true
		}
	}

	// Si menciona función
	if containsAny(step, "función", "function") {
		fn := g.extractFunctionName(step)
		if fn != "" {
			return &GroundedToolCall{
				ToolName: "find_function",
				Args: map[string]interface{}{
					"name": fn,
				},
			}, true
		}
	}

	// Búsqueda global
	if containsAny(step, "buscar", "search") {
		return &GroundedToolCall{
			ToolName: "search_in_project",
			Args: map[string]interface{}{
				"query": step,
			},
		}, true
	}

	// Fallback semántico
	return g.semanticFallback(step)
}

//
// ============================================================
//  SEMANTIC FALLBACK
// ============================================================
//

func (g *ContextualToolGrounder) semanticFallback(step string) (*GroundedToolCall, bool) {
	s := strings.ToLower(step)

	if containsAny(s, "resumir", "summary") {
		return &GroundedToolCall{
			ToolName: "summarize_text",
			Args: map[string]interface{}{
				"text": step,
			},
		}, true
	}

	if containsAny(s, "explicar", "explain") {
		return &GroundedToolCall{
			ToolName: "explain_code",
			Args: map[string]interface{}{
				"input": step,
			},
		}, true
	}

	return nil, false
}

//
// ============================================================
//  HELPERS
// ============================================================
//

func (g *ContextualToolGrounder) extractFileName(step string) string {
	for _, w := range strings.Fields(step) {
		if strings.Contains(w, ".go") || strings.Contains(w, ".md") {
			return w
		}
	}
	return ""
}

func (g *ContextualToolGrounder) extractFunctionName(step string) string {
	for _, w := range strings.Fields(step) {
		if strings.HasPrefix(w, "func") {
			return w
		}
	}
	return ""
}

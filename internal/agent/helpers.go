package agent

import (
	"opencode-lite/internal/tools"
	"strings"
)

func extractPattern(goal string) string {
	// Ejemplo simple: buscar "algo"
	parts := strings.Split(goal, "\"")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ".*"
}

func extractFile(goal string) string {
	// Busca algo como archivo.go
	words := strings.Fields(goal)
	for _, w := range words {
		if strings.HasSuffix(w, ".go") {
			return w
		}
	}
	return ""
}

func extractFiles(goal string) []string {
	// Extrae múltiples archivos
	files := []string{}
	words := strings.Fields(goal)
	for _, w := range words {
		if strings.HasSuffix(w, ".go") {
			files = append(files, w)
		}
	}
	return files
}

func extractRename(goal string) (string, string) {
	// renombrar foo a bar
	words := strings.Fields(goal)
	var old, new string
	for i, w := range words {
		if w == "a" || w == "to" {
			if i > 0 {
				old = words[i-1]
			}
			if i < len(words)-1 {
				new = words[i+1]
			}
		}
	}
	return old, new
}

func containsAny(s string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}

func extractNewPath(goal string) string {
	// mover archivo foo.go a /nuevo/lugar
	words := strings.Fields(goal)
	for i, w := range words {
		if w == "a" || w == "to" {
			if i < len(words)-1 {
				return words[i+1]
			}
		}
	}
	return ""
}

func (rt *AgentRuntime) executeTool(name string, args map[string]interface{}, ctx *AgentContext) tools.ToolResult {
	toolFn := tools.ToolRegistry[name]
	result := toolFn(args)

	ctx.History = append(ctx.History, AgentStep{
		Thought: "Ejecutando " + name,
		Action:  name,
		Input:   args,
		Output:  result,
	})

	ctx.LastResult = result
	return result
}

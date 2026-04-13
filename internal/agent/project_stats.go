package agent

import (
	"opencode-lite/internal/tools"
	"strings"
)

type ProjectStats struct {
	FileCount    int
	PackageCount int
	LargeFiles   int
}

func AnalyzeProjectSize() ProjectStats {
	stats := ProjectStats{}

	// 1. Obtener file tree real
	treeTool, ok := tools.ToolRegistry["file_tree"]
	if !ok {
		return ProjectStats{FileCount: 0, PackageCount: 0, LargeFiles: 0}
	}

	result := treeTool(map[string]interface{}{
		"root": "workspace",
	})

	if result.Error != "" {
		return stats
	}

	// === result.Result es un map[string]interface{} ===
	payload, ok := result.Result.(map[string]interface{})
	if !ok {
		return stats
	}

	// === extraemos el campo "tree" ===
	rawTree, ok := payload["tree"]
	if !ok {
		return stats
	}

	// === convertir a []string ===
	var files []string

	// Caso 1: ya es []string
	if arr, ok := rawTree.([]string); ok {
		files = arr

		// Caso 2: viene como []interface{} (muy común)
	} else if arr, ok := rawTree.([]interface{}); ok {
		files = make([]string, len(arr))
		for i, v := range arr {
			files[i] = v.(string)
		}

	} else {
		return stats
	}

	stats.FileCount = len(files)

	// 2. Contar paquetes Go
	packages := map[string]bool{}
	for _, f := range files {
		if strings.HasSuffix(f, ".go") {
			dir := extractDir(f)
			packages[dir] = true
		}
	}
	stats.PackageCount = len(packages)

	// 3. Detectar archivos grandes (heurística)
	large := 0
	for _, f := range files {
		if strings.HasSuffix(f, ".go") {
			if isLargeFile(f) {
				large++
			}
		}
	}
	stats.LargeFiles = large

	return stats
}

func extractDir(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) <= 1 {
		return "."
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

// Heurística: contar líneas usando search_regex_multi
func isLargeFile(path string) bool {
	tool, ok := tools.ToolRegistry["search_regex_multi"]
	if !ok {
		return false
	}

	result := tool(map[string]interface{}{
		"path":    path,
		"pattern": "\n",
	})

	if result.Error != "" {
		return false
	}

	count, ok := result.Result.(int)
	if !ok {
		return false
	}

	return count > 300
}

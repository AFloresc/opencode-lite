package agent

import (
	"strings"
)

type AgentPolicy interface {
	Decide(ctx *AgentContext) (toolName string, args map[string]interface{}, done bool)
}

type RuleBasedPolicy struct{}

func (p RuleBasedPolicy) Decide(ctx *AgentContext) (string, map[string]interface{}, bool) {
	goal := strings.ToLower(ctx.Goal)

	// ============================
	// 🔍 BÚSQUEDA Y ANÁLISIS DE TEXTO
	// ============================
	if containsAny(goal, "buscar", "search", "encontrar", "grep") {
		pattern := extractPattern(goal)
		return "search_regex_multi", map[string]interface{}{
			"path":    "workspace",
			"pattern": pattern,
		}, false
	}

	if containsAny(goal, "contar funciones", "count funcs") {
		file := extractFile(goal)
		return "count_funcs", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "contar imports", "count imports") {
		file := extractFile(goal)
		return "count_imports", map[string]interface{}{
			"path": file,
		}, false
	}

	// ============================
	// 📊 MÉTRICAS Y ANÁLISIS
	// ============================
	if containsAny(goal, "métricas", "metrics", "estadísticas", "stats") {
		file := extractFile(goal)
		return "analysis_metrics", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "dependencias", "dependencies") {
		return "analysis_dependencies", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	if containsAny(goal, "complejidad", "cyclomatic") {
		file := extractFile(goal)
		return "analysis_cyclomatic", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "dead code", "código muerto") {
		return "analysis_dead_code", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	// ============================
	// 📁 ARCHIVOS Y PROYECTO
	// ============================
	if containsAny(goal, "listar archivos", "file tree", "archivos", "listar") {
		return "file_tree", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	if containsAny(goal, "proyecto", "analizar proyecto", "project analysis") {
		return "project_stats", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	if containsAny(goal, "archivos grandes", "largest files") {
		return "largest_files", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	// ============================
	// 🛠️ REFACTOR
	// ============================
	if containsAny(goal, "renombrar", "rename") {
		symbol, newName := extractRename(goal)
		return "refactor_rename_symbol", map[string]interface{}{
			"symbol":   symbol,
			"new_name": newName,
		}, false
	}

	if containsAny(goal, "mover archivo", "move file") {
		file := extractFile(goal)
		newPath := extractNewPath(goal)
		return "refactor_move_file", map[string]interface{}{
			"path":     file,
			"new_path": newPath,
		}, false
	}

	if containsAny(goal, "split", "dividir archivo") {
		file := extractFile(goal)
		return "refactor_split_file", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "merge", "fusionar archivos") {
		files := extractFiles(goal)
		return "refactor_merge_files", map[string]interface{}{
			"files": files,
		}, false
	}

	// ============================
	// 🧠 SEMÁNTICA
	// ============================
	if containsAny(goal, "extraer funciones", "extract functions") {
		file := extractFile(goal)
		return "extract_functions", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "extraer tipos", "extract types") {
		file := extractFile(goal)
		return "extract_types", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "extraer comentarios", "extract comments") {
		file := extractFile(goal)
		return "extract_comments_block", map[string]interface{}{
			"path": file,
		}, false
	}

	if containsAny(goal, "indexar", "semantic index") {
		return "semantic_index", map[string]interface{}{
			"root": "workspace",
		}, false
	}

	// ============================
	// ❌ SIN REGLA → TERMINAR
	// ============================
	return "", nil, true
}

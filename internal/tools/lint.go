package tools

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// lintCodeTool
// Linter básico según extensión:
// - .go     → lint sintáctico simple
// - .json   → validación + estructura
// - .yaml   → validación mínima
// - otros   → chequeos genéricos
// ------------------------------------------------------------
func lintCodeTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"lint_code", nil, "falta argumento obligatorio: path"}
	}

	path := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"lint_code", nil, err.Error()}
	}

	ext := strings.ToLower(filepath.Ext(path))

	var issues []string

	switch ext {
	case ".go":
		issues = lintGo(content)
	case ".json":
		issues = lintJSON(content)
	case ".yaml", ".yml":
		issues = lintYAML(content)
	default:
		issues = lintGeneric(content)
	}

	return ToolResult{"lint_code", map[string]interface{}{
		"path":   path,
		"issues": issues,
		"count":  len(issues),
	}, ""}
}

//
// ------------------------------------------------------------
// IMPLEMENTACIONES INTERNAS
// ------------------------------------------------------------
//

// ------------------------------------------------------------
// Lint para Go (simple, sin dependencias externas)
// ------------------------------------------------------------
func lintGo(src string) []string {
	var issues []string
	lines := strings.Split(src, "\n")

	// 1. Líneas demasiado largas
	for i, l := range lines {
		if len(l) > 120 {
			issues = append(issues,
				fmtIssue(i, "línea demasiado larga (>120 caracteres)"))
		}
	}

	// 2. Tabs vs spaces
	for i, l := range lines {
		if strings.Contains(l, "\t") {
			issues = append(issues,
				fmtIssue(i, "usa tabulaciones; Go recomienda tabs pero mezcla puede ser inconsistente"))
		}
	}

	// 3. Comentarios TODO
	for i, l := range lines {
		if strings.Contains(strings.ToUpper(l), "TODO") {
			issues = append(issues,
				fmtIssue(i, "TODO encontrado"))
		}
	}

	// 4. Funciones sin comentarios (heurística)
	reFunc := regexp.MustCompile(`^func\s+([A-Za-z0-9_]+)`)
	for i, l := range lines {
		if reFunc.MatchString(l) {
			if i == 0 || !strings.HasPrefix(strings.TrimSpace(lines[i-1]), "//") {
				issues = append(issues,
					fmtIssue(i, "función sin comentario previo"))
			}
		}
	}

	return issues
}

// ------------------------------------------------------------
// Lint para JSON
// ------------------------------------------------------------
func lintJSON(src string) []string {
	var issues []string

	var obj interface{}
	err := json.Unmarshal([]byte(src), &obj)
	if err != nil {
		issues = append(issues, "JSON inválido: "+err.Error())
		return issues
	}

	// JSON válido → sugerencias
	if !strings.HasPrefix(strings.TrimSpace(src), "{") {
		issues = append(issues, "JSON debería comenzar con '{'")
	}

	return issues
}

// ------------------------------------------------------------
// Lint para YAML (validación mínima)
// ------------------------------------------------------------
func lintYAML(src string) []string {
	var issues []string
	lines := strings.Split(src, "\n")

	// 1. Indentación inconsistente
	for i, l := range lines {
		if strings.HasPrefix(l, " ") && strings.Contains(l, "\t") {
			issues = append(issues,
				fmtIssue(i, "mezcla de tabs y espacios en YAML"))
		}
	}

	// 2. Claves sin valor
	reKey := regexp.MustCompile(`^[A-Za-z0-9_-]+:\s*$`)
	for i, l := range lines {
		if reKey.MatchString(l) {
			issues = append(issues,
				fmtIssue(i, "clave YAML sin valor"))
		}
	}

	return issues
}

// ------------------------------------------------------------
// Lint genérico (texto)
// ------------------------------------------------------------
func lintGeneric(src string) []string {
	var issues []string
	lines := strings.Split(src, "\n")

	// 1. Espacios al final
	for i, l := range lines {
		if strings.HasSuffix(l, " ") {
			issues = append(issues,
				fmtIssue(i, "espacios al final de línea"))
		}
	}

	// 2. Líneas vacías duplicadas
	lastEmpty := false
	for i, l := range lines {
		if strings.TrimSpace(l) == "" {
			if lastEmpty {
				issues = append(issues,
					fmtIssue(i, "línea vacía duplicada"))
			}
			lastEmpty = true
		} else {
			lastEmpty = false
		}
	}

	return issues
}

// ------------------------------------------------------------
// Helper para formatear issues
// ------------------------------------------------------------
func fmtIssue(line int, msg string) string {
	return fmt.Sprintf("línea %d: %s", line+1, msg)
}

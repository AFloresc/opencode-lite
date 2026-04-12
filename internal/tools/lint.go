package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// lintCodeTool

// lint_code: analiza un archivo y devuelve advertencias comunes
func lintCodeTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"lint_code", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"lint_code", nil, "el argumento 'path' debe ser string"}
	}

	lang := ""
	if langRaw, ok := args["lang"]; ok {
		lang, _ = langRaw.(string)
	}

	fullPath := filepath.Join("workspace", path)

	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"lint_code", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"lint_code", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")

	// Detectar lenguaje por extensión si no se especifica
	if lang == "" {
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go":
			lang = "go"
		case ".json":
			lang = "json"
		case ".yaml", ".yml":
			lang = "yaml"
		default:
			lang = "generic"
		}
	}

	var warnings []map[string]interface{}

	// ---------------------------------------------------------
	// Reglas genéricas
	// ---------------------------------------------------------
	for i, line := range lines {
		lineno := i + 1

		// Espacios al final
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			warnings = append(warnings, map[string]interface{}{
				"line":    lineno,
				"type":    "trailing_whitespace",
				"message": "La línea tiene espacios al final",
			})
		}

		// Líneas demasiado largas
		if len(line) > 120 {
			warnings = append(warnings, map[string]interface{}{
				"line":    lineno,
				"type":    "line_too_long",
				"message": "La línea supera los 120 caracteres",
			})
		}

		// TODOs
		if strings.Contains(line, "TODO") {
			warnings = append(warnings, map[string]interface{}{
				"line":    lineno,
				"type":    "todo",
				"message": "Hay un TODO pendiente",
			})
		}
	}

	// ---------------------------------------------------------
	// Reglas específicas para Go
	// ---------------------------------------------------------
	if lang == "go" {
		// Import no usado (heurístico)
		importRe := regexp.MustCompile(`"([^"]+)"`)
		imports := importRe.FindAllStringSubmatch(content, -1)

		for _, imp := range imports {
			pkg := filepath.Base(imp[1])
			if !strings.Contains(content, pkg+".") {
				warnings = append(warnings, map[string]interface{}{
					"type":    "unused_import",
					"message": fmt.Sprintf("Import '%s' posiblemente no usado", imp[1]),
				})
			}
		}

		// Funciones no usadas (heurístico)
		funcRe := regexp.MustCompile(`func\s+([A-Za-z0-9_]+)\s*\(`)
		funcs := funcRe.FindAllStringSubmatch(content, -1)

		for _, f := range funcs {
			name := f[1]
			if !strings.Contains(content, name+"(") || strings.Count(content, name+"(") == 1 {
				warnings = append(warnings, map[string]interface{}{
					"type":    "unused_function",
					"message": fmt.Sprintf("La función '%s' podría no estar siendo usada", name),
				})
			}
		}
	}

	// ---------------------------------------------------------
	// Reglas JSON
	// ---------------------------------------------------------
	if lang == "json" {
		var js interface{}
		if err := json.Unmarshal(contentBytes, &js); err != nil {
			warnings = append(warnings, map[string]interface{}{
				"type":    "invalid_json",
				"message": fmt.Sprintf("JSON inválido: %v", err),
			})
		}
	}

	// ---------------------------------------------------------
	// Reglas YAML (básicas)
	// ---------------------------------------------------------
	if lang == "yaml" {
		if strings.Contains(content, "\t") {
			warnings = append(warnings, map[string]interface{}{
				"type":    "yaml_tabs",
				"message": "YAML no debe usar tabs para indentación",
			})
		}
	}

	return ToolResult{
		ToolName: "lint_code",
		Result: map[string]interface{}{
			"path":     path,
			"lang":     lang,
			"warnings": warnings,
			"count":    len(warnings),
		},
	}
}

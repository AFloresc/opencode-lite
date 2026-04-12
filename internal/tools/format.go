package tools

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

// formatCodeTool
// format_code: autoformatea código según el lenguaje detectado o especificado
func formatCodeTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"format_code", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"format_code", nil, "el argumento 'path' debe ser string"}
	}

	lang := ""
	if langRaw, ok := args["lang"]; ok {
		lang, _ = langRaw.(string)
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo
	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"format_code", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"format_code", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)

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

	var formatted string

	switch lang {

	// ---------------------------------------------------------
	// GO FORMATTER (real gofmt)
	// ---------------------------------------------------------
	case "go":
		formattedBytes, err := format.Source([]byte(content))
		if err != nil {
			return ToolResult{"format_code", nil, fmt.Sprintf("error formateando Go: %v", err)}
		}
		formatted = string(formattedBytes)

	// ---------------------------------------------------------
	// JSON FORMATTER
	// ---------------------------------------------------------
	case "json":
		var obj interface{}
		if err := json.Unmarshal([]byte(content), &obj); err != nil {
			return ToolResult{"format_code", nil, fmt.Sprintf("JSON inválido: %v", err)}
		}
		out, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return ToolResult{"format_code", nil, fmt.Sprintf("error formateando JSON: %v", err)}
		}
		formatted = string(out)

	// ---------------------------------------------------------
	// YAML FORMATTER (simple indent)
	// ---------------------------------------------------------
	case "yaml":
		// YAML no tiene formateador estándar en Go sin dependencias externas.
		// Hacemos un trim + normalización básica.
		lines := strings.Split(content, "\n")
		var cleaned []string
		for _, l := range lines {
			cleaned = append(cleaned, strings.TrimRight(l, " \t"))
		}
		formatted = strings.Join(cleaned, "\n")

	// ---------------------------------------------------------
	// GENERIC FORMATTER
	// ---------------------------------------------------------
	default:
		lines := strings.Split(content, "\n")
		var cleaned []string
		for _, l := range lines {
			cleaned = append(cleaned, strings.TrimRight(l, " \t"))
		}
		formatted = strings.Join(cleaned, "\n")
	}

	// Guardar archivo formateado
	err = os.WriteFile(fullPath, []byte(formatted), 0644)
	if err != nil {
		return ToolResult{"format_code", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "format_code",
		Result:   fmt.Sprintf("archivo '%s' formateado correctamente como '%s'", path, lang),
	}
}

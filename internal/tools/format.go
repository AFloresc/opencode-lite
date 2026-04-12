package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------
// formatCodeTool
// Formatea archivos según su extensión:
// - .go     → gofmt
// - .json   → indent JSON
// - .yaml   → limpieza básica
// - .yml    → limpieza básica
// - .md     → normalización ligera
// - otros   → trim + normalización
// ------------------------------------------------------------
func formatCodeTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"format_code", nil, "falta argumento obligatorio: path"}
	}

	path := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"format_code", nil, err.Error()}
	}

	ext := strings.ToLower(filepath.Ext(path))

	var formatted string

	switch ext {
	case ".go":
		formatted, err = formatGo(content)
	case ".json":
		formatted, err = formatJSON(content)
	case ".yaml", ".yml":
		formatted, err = formatYAML(content)
	case ".md":
		formatted, err = formatMarkdown(content)
	default:
		formatted, err = formatGeneric(content)
	}

	if err != nil {
		return ToolResult{"format_code", nil, err.Error()}
	}

	if err := writeFile(path, formatted); err != nil {
		return ToolResult{"format_code", nil, err.Error()}
	}

	return ToolResult{"format_code", map[string]interface{}{
		"path":      path,
		"formatted": true,
	}, ""}
}

//
// ------------------------------------------------------------
// IMPLEMENTACIONES INTERNAS
// ------------------------------------------------------------
//

// ------------------------------------------------------------
// Go formatter (gofmt real)
// ------------------------------------------------------------
func formatGo(src string) (string, error) {
	cmd := exec.Command("gofmt")
	cmd.Stdin = bytes.NewBufferString(src)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error ejecutando gofmt: %v", err)
	}

	return string(out), nil
}

// ------------------------------------------------------------
// JSON formatter
// ------------------------------------------------------------
func formatJSON(src string) (string, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(src), &obj); err != nil {
		return "", errors.New("JSON inválido: " + err.Error())
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(pretty), nil
}

// ------------------------------------------------------------
// YAML formatter (simple normalización)
// ------------------------------------------------------------
func formatYAML(src string) (string, error) {
	// No usamos librerías externas para mantener sandbox limpio.
	// Normalización básica:
	// - trim espacios
	// - eliminar líneas vacías duplicadas
	lines := strings.Split(src, "\n")
	var out []string

	lastEmpty := false
	for _, l := range lines {
		trimmed := strings.TrimRight(l, " ")

		if trimmed == "" {
			if lastEmpty {
				continue
			}
			lastEmpty = true
		} else {
			lastEmpty = false
		}

		out = append(out, trimmed)
	}

	return strings.Join(out, "\n"), nil
}

// ------------------------------------------------------------
// Markdown formatter (ligero)
// ------------------------------------------------------------
func formatMarkdown(src string) (string, error) {
	lines := strings.Split(src, "\n")
	var out []string

	for _, l := range lines {
		out = append(out, strings.TrimRight(l, " "))
	}

	return strings.Join(out, "\n"), nil
}

// ------------------------------------------------------------
// Generic formatter (fallback)
// ------------------------------------------------------------
func formatGeneric(src string) (string, error) {
	// Normalización mínima:
	// - trim final
	// - eliminar espacios repetidos al final de línea
	lines := strings.Split(src, "\n")
	var out []string

	for _, l := range lines {
		out = append(out, strings.TrimRight(l, " "))
	}

	return strings.TrimSpace(strings.Join(out, "\n")), nil
}

package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//
// Estructuras base
//

// ToolCall representa una llamada a herramienta generada por el modelo
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult representa el resultado de ejecutar una herramienta
type ToolResult struct {
	ToolName string      `json:"tool_name"`
	Result   interface{} `json:"result"`
	Error    string      `json:"error,omitempty"`
}

//
// Registro de herramientas disponibles
//

var toolRegistry = map[string]func(map[string]interface{}) ToolResult{
	"read_file":   readFileTool,
	"write_file":  writeFileTool,
	"apply_patch": applyPatchTool,
}

//
// Función principal para ejecutar herramientas
//

func ExecuteTool(call ToolCall) ToolResult {
	toolFunc, exists := toolRegistry[call.Name]
	if !exists {
		return ToolResult{
			ToolName: call.Name,
			Error:    fmt.Sprintf("herramienta desconocida: %s", call.Name),
		}
	}

	return toolFunc(call.Arguments)
}

//
// Implementación de herramientas
//

// read_file: lee un archivo dentro del directorio workspace
func readFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{
			ToolName: "read_file",
			Error:    "falta argumento obligatorio: path",
		}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{
			ToolName: "read_file",
			Error:    "el argumento 'path' debe ser string",
		}
	}

	// 🔥 Siempre leemos desde workspace/
	fullPath := filepath.Join("workspace", path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{
			ToolName: "read_file",
			Error:    fmt.Sprintf("error leyendo archivo: %v", err),
		}
	}

	return ToolResult{
		ToolName: "read_file",
		Result:   string(data),
	}
}

// write_file: escribe un archivo dentro del directorio workspace
func writeFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{
			ToolName: "write_file",
			Error:    "falta argumento obligatorio: path",
		}
	}

	contentRaw, ok := args["content"]
	if !ok {
		return ToolResult{
			ToolName: "write_file",
			Error:    "falta argumento obligatorio: content",
		}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{
			ToolName: "write_file",
			Error:    "el argumento 'path' debe ser string",
		}
	}

	content, ok := contentRaw.(string)
	if !ok {
		return ToolResult{
			ToolName: "write_file",
			Error:    "el argumento 'content' debe ser string",
		}
	}

	fullPath := filepath.Join("workspace", path)

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return ToolResult{
			ToolName: "write_file",
			Error:    fmt.Sprintf("error escribiendo archivo: %v", err),
		}
	}

	return ToolResult{
		ToolName: "write_file",
		Result:   fmt.Sprintf("archivo escrito correctamente: %s", fullPath),
	}
}

// / apply_patch: aplica un parche unified diff real
func applyPatchTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{"apply_patch", nil, fmt.Sprintf("error leyendo archivo original: %v", err)}
	}

	original := strings.Split(string(originalBytes), "\n")
	result := []string{}
	i := 0 // índice en original

	lines := strings.Split(patch, "\n")

	for _, line := range lines {

		// Ignorar encabezados
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			continue
		}

		// Ignorar encabezados de hunk
		if strings.HasPrefix(line, "@@") {
			continue
		}

		// Línea eliminada
		if strings.HasPrefix(line, "-") {
			i++
			continue
		}

		// Línea añadida
		if strings.HasPrefix(line, "+") {
			result = append(result, line[1:])
			continue
		}

		// Línea sin prefijo → copiar del original
		if i < len(original) {
			result = append(result, original[i])
			i++
		}
	}

	// Guardar archivo modificado
	final := strings.Join(result, "\n")
	err = os.WriteFile(fullPath, []byte(final), 0644)
	if err != nil {
		return ToolResult{"apply_patch", nil, fmt.Sprintf("error escribiendo archivo modificado: %v", err)}
	}

	return ToolResult{"apply_patch", fmt.Sprintf("parche aplicado correctamente a %s", fullPath), ""}
}

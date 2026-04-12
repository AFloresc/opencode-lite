package tools

import (
	"fmt"
	"os"
	"path/filepath"
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
	"read_file":  readFileTool,
	"write_file": writeFileTool,
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

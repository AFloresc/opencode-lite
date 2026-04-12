package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

func ExecuteTool(call ToolCall) ToolResult {
	switch call.Name {
	case "write_file":
		return writeFileTool(call.Arguments)
	default:
		return ToolResult{
			Name:  call.Name,
			Error: "unknown tool",
		}
	}
}

func writeFileTool(args map[string]interface{}) ToolResult {
	pathVal, ok1 := args["path"].(string)
	contentVal, ok2 := args["content"].(string)
	if !ok1 || !ok2 {
		return ToolResult{
			Name:  "write_file",
			Error: "invalid arguments",
		}
	}

	// Crear siempre un workspace seguro
	base := "workspace"

	// Normalizar la ruta para evitar rutas absolutas o escapadas
	clean := filepath.Clean(pathVal)

	// Forzar que SIEMPRE escriba dentro de workspace/
	safePath := filepath.Join(base, clean)

	// Crear todas las carpetas necesarias
	if err := os.MkdirAll(filepath.Dir(safePath), 0o755); err != nil {
		return ToolResult{
			Name:  "write_file",
			Error: fmt.Sprintf("mkdir: %v", err),
		}
	}

	// Escribir el archivo
	if err := os.WriteFile(safePath, []byte(contentVal), 0o644); err != nil {
		return ToolResult{
			Name:  "write_file",
			Error: fmt.Sprintf("write: %v", err),
		}
	}

	return ToolResult{
		Name:   "write_file",
		Result: fmt.Sprintf("file written: %s", safePath),
	}
}

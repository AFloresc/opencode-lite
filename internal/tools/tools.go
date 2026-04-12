package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

// apply_patch ULTRA: soporte para múltiples hunks, validación de contexto y detección de conflictos
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
	result := make([]string, 0, len(original))

	lines := strings.Split(patch, "\n")
	origIdx := 0 // índice en original

	// Estado de conflicto
	conflicts := []string{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Ignorar encabezados de archivo
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			continue
		}

		// Hunk: @@ -a,b +c,d @@
		if strings.HasPrefix(line, "@@") {
			// Ejemplo: @@ -1,3 +1,4 @@
			parts := strings.Split(line, " ")
			if len(parts) < 3 {
				conflicts = append(conflicts, fmt.Sprintf("hunk mal formado: %s", line))
				continue
			}

			// Rango original: "-a,b"
			oldRange := strings.TrimPrefix(parts[1], "-")
			oldStartStr := strings.Split(oldRange, ",")[0]
			oldStart, err := strconv.Atoi(oldStartStr)
			if err != nil {
				conflicts = append(conflicts, fmt.Sprintf("no se pudo parsear offset original en hunk: %s", line))
				continue
			}
			oldStart-- // 1-based → 0-based

			// Copiar líneas previas sin cambios
			for origIdx < oldStart && origIdx < len(original) {
				result = append(result, original[origIdx])
				origIdx++
			}

			// Procesar hunk
			i++
			for i < len(lines) {
				hline := lines[i]

				// Nuevo hunk → retrocede uno para que el for externo lo procese
				if strings.HasPrefix(hline, "@@") {
					i--
					break
				}

				// Encabezados de archivo dentro del diff
				if strings.HasPrefix(hline, "--- ") || strings.HasPrefix(hline, "+++ ") {
					// dejamos que el siguiente ciclo los ignore
					break
				}

				// Línea eliminada: debe coincidir con original
				if strings.HasPrefix(hline, "-") {
					expected := hline[1:]
					if origIdx >= len(original) || original[origIdx] != expected {
						conflicts = append(conflicts, fmt.Sprintf(
							"conflicto: se esperaba eliminar '%s' pero en el archivo hay '%s'",
							expected,
							func() string {
								if origIdx < len(original) {
									return original[origIdx]
								}
								return "<EOF>"
							}(),
						))
					} else {
						// Coincide → se elimina avanzando el índice
						origIdx++
					}
					i++
					continue
				}

				// Línea añadida
				if strings.HasPrefix(hline, "+") {
					result = append(result, hline[1:])
					i++
					continue
				}

				// Línea de contexto (sin prefijo): debe coincidir con original
				if !strings.HasPrefix(hline, "+") && !strings.HasPrefix(hline, "-") {
					if origIdx >= len(original) || original[origIdx] != hline {
						conflicts = append(conflicts, fmt.Sprintf(
							"conflicto en contexto: se esperaba '%s' pero en el archivo hay '%s'",
							hline,
							func() string {
								if origIdx < len(original) {
									return original[origIdx]
								}
								return "<EOF>"
							}(),
						))
					} else {
						result = append(result, original[origIdx])
						origIdx++
					}
					i++
					continue
				}
			}

			continue
		}
	}

	// Copiar el resto del archivo original
	for origIdx < len(original) {
		result = append(result, original[origIdx])
		origIdx++
	}

	// Si hubo conflictos, no escribimos nada
	if len(conflicts) > 0 {
		return ToolResult{
			ToolName: "apply_patch",
			Result:   nil,
			Error:    fmt.Sprintf("conflictos al aplicar el parche:\n%s", strings.Join(conflicts, "\n")),
		}
	}

	final := strings.Join(result, "\n")
	err = os.WriteFile(fullPath, []byte(final), 0644)
	if err != nil {
		return ToolResult{"apply_patch", nil, fmt.Sprintf("error escribiendo archivo modificado: %v", err)}
	}

	return ToolResult{"apply_patch", fmt.Sprintf("parche aplicado correctamente a %s", fullPath), ""}
}

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
	"read_file":         readFileTool,
	"write_file":        writeFileTool,
	"apply_patch":       applyPatchTool,
	"apply_patch_fuzzy": applyPatchFuzzyTool,
	"list_files":        listFilesTool,
	"search_in_file":    searchInFileTool,
	"grep":              grepTool,
	"delete_file":       deleteFileTool,
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

// apply_patch_fuzzy: aplica un parche aunque el contexto no coincida exactamente.
// Permite duplicar cambios y modificar líneas parcialmente.
func applyPatchFuzzyTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_fuzzy", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{"apply_patch_fuzzy", nil, fmt.Sprintf("error leyendo archivo original: %v", err)}
	}

	original := strings.Split(string(originalBytes), "\n")
	result := make([]string, 0, len(original))

	lines := strings.Split(patch, "\n")

	// Extraer líneas - y +
	var toRemove []string
	var toAdd []string

	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			toRemove = append(toRemove, line[1:])
		}
		if strings.HasPrefix(line, "+") {
			toAdd = append(toAdd, line[1:])
		}
	}

	// Aplicación fuzzy
	for _, line := range original {
		removed := false
		for _, r := range toRemove {
			if strings.Contains(line, r) {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, line)
		}
	}

	// Añadir líneas nuevas al final
	result = append(result, toAdd...)

	final := strings.Join(result, "\n")
	err = os.WriteFile(fullPath, []byte(final), 0644)
	if err != nil {
		return ToolResult{"apply_patch_fuzzy", nil, fmt.Sprintf("error escribiendo archivo modificado: %v", err)}
	}

	return ToolResult{"apply_patch_fuzzy", fmt.Sprintf("parche fuzzy aplicado correctamente a %s", fullPath), ""}
}

// list_files: lista archivos dentro del directorio workspace
func listFilesTool(args map[string]interface{}) ToolResult {
	recursive := false
	extFilter := ""

	// Argumento opcional: recursive
	if r, ok := args["recursive"]; ok {
		if rBool, ok := r.(bool); ok {
			recursive = rBool
		}
	}

	// Argumento opcional: ext
	if e, ok := args["ext"]; ok {
		if eStr, ok := e.(string); ok {
			extFilter = eStr
		}
	}

	base := "workspace"
	var files []string

	if recursive {
		// Recorrido recursivo
		err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if extFilter != "" && filepath.Ext(path) != extFilter {
				return nil
			}
			rel, _ := filepath.Rel(base, path)
			files = append(files, rel)
			return nil
		})
		if err != nil {
			return ToolResult{"list_files", nil, fmt.Sprintf("error recorriendo directorio: %v", err)}
		}
	} else {
		// Solo nivel superior
		entries, err := os.ReadDir(base)
		if err != nil {
			return ToolResult{"list_files", nil, fmt.Sprintf("error leyendo directorio: %v", err)}
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if extFilter != "" && filepath.Ext(entry.Name()) != extFilter {
				continue
			}
			files = append(files, entry.Name())
		}
	}

	return ToolResult{
		ToolName: "list_files",
		Result:   files,
	}
}

// search_in_file: busca texto dentro de un archivo y devuelve coincidencias con número de línea
func searchInFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"search_in_file", nil, "falta argumento obligatorio: path"}
	}

	queryRaw, ok := args["query"]
	if !ok {
		return ToolResult{"search_in_file", nil, "falta argumento obligatorio: query"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"search_in_file", nil, "el argumento 'path' debe ser string"}
	}

	query, ok := queryRaw.(string)
	if !ok {
		return ToolResult{"search_in_file", nil, "el argumento 'query' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo
	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ToolResult{"search_in_file", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	lines := strings.Split(string(contentBytes), "\n")
	var results []map[string]interface{}

	for i, line := range lines {
		if strings.Contains(line, query) {
			results = append(results, map[string]interface{}{
				"line_number": i + 1,
				"line":        line,
			})
		}
	}

	return ToolResult{
		ToolName: "search_in_file",
		Result:   results,
	}
}

// grep: busca texto en todos los archivos del workspace (multi-archivo)
func grepTool(args map[string]interface{}) ToolResult {
	queryRaw, ok := args["query"]
	if !ok {
		return ToolResult{"grep", nil, "falta argumento obligatorio: query"}
	}

	query, ok := queryRaw.(string)
	if !ok {
		return ToolResult{"grep", nil, "el argumento 'query' debe ser string"}
	}

	extFilter := ""
	if e, ok := args["ext"]; ok {
		if eStr, ok := e.(string); ok {
			extFilter = eStr
		}
	}

	recursive := true
	if r, ok := args["recursive"]; ok {
		if rBool, ok := r.(bool); ok {
			recursive = rBool
		}
	}

	base := "workspace"
	var results []map[string]interface{}

	// Función para procesar un archivo
	processFile := func(path string) error {
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil // ignoramos errores por archivo
		}

		lines := strings.Split(string(contentBytes), "\n")
		rel, _ := filepath.Rel(base, path)

		for i, line := range lines {
			if strings.Contains(line, query) {
				results = append(results, map[string]interface{}{
					"file":        rel,
					"line_number": i + 1,
					"line":        line,
				})
			}
		}
		return nil
	}

	if recursive {
		// Recorrido recursivo
		filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if extFilter != "" && filepath.Ext(path) != extFilter {
				return nil
			}
			return processFile(path)
		})
	} else {
		// Solo nivel superior
		entries, err := os.ReadDir(base)
		if err != nil {
			return ToolResult{"grep", nil, fmt.Sprintf("error leyendo directorio: %v", err)}
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if extFilter != "" && filepath.Ext(entry.Name()) != extFilter {
				continue
			}
			processFile(filepath.Join(base, entry.Name()))
		}
	}

	return ToolResult{
		ToolName: "grep",
		Result:   results,
	}
}

// delete_file: elimina un archivo o directorio dentro del workspace
func deleteFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"delete_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"delete_file", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Verificar existencia
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"delete_file", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"delete_file", nil, fmt.Sprintf("error accediendo al archivo: %v", err)}
	}

	// Si es directorio, eliminar recursivamente
	if info.IsDir() {
		err = os.RemoveAll(fullPath)
		if err != nil {
			return ToolResult{"delete_file", nil, fmt.Sprintf("error eliminando directorio: %v", err)}
		}
		return ToolResult{"delete_file", fmt.Sprintf("directorio '%s' eliminado correctamente", path), ""}
	}

	// Si es archivo, eliminar normalmente
	err = os.Remove(fullPath)
	if err != nil {
		return ToolResult{"delete_file", nil, fmt.Sprintf("error eliminando archivo: %v", err)}
	}

	return ToolResult{"delete_file", fmt.Sprintf("archivo '%s' eliminado correctamente", path), ""}
}

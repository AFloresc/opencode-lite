package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	"read_file":              readFileTool,
	"write_file":             writeFileTool,
	"apply_patch":            applyPatchTool,
	"apply_patch_fuzzy":      applyPatchFuzzyTool,
	"apply_patch_auto":       applyPatchAutoTool,
	"list_files":             listFilesTool,
	"search_in_file":         searchInFileTool,
	"grep":                   grepTool,
	"delete_file":            deleteFileTool,
	"rename_file":            renameFileTool,
	"copy_file":              copyFileTool,
	"move_file":              moveFileTool,
	"create_file":            createFileTool,
	"file_exists":            fileExistsTool,
	"read_dir":               readDirTool,
	"append_file":            appendFileTool,
	"truncate_file":          truncateFileTool,
	"stat_file":              statFileTool,
	"touch_file":             touchFileTool,
	"search_replace":         searchReplaceTool,
	"diff_files":             diffFilesTool,
	"zip_dir":                zipDirTool,
	"unzip":                  unzipTool,
	"search_regex":           searchRegexTool,
	"apply_patch_structured": applyPatchStructuredTool,
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

// rename_file: renombra un archivo o directorio dentro del workspace
func renameFileTool(args map[string]interface{}) ToolResult {
	fromRaw, ok := args["from"]
	if !ok {
		return ToolResult{"rename_file", nil, "falta argumento obligatorio: from"}
	}

	toRaw, ok := args["to"]
	if !ok {
		return ToolResult{"rename_file", nil, "falta argumento obligatorio: to"}
	}

	from, ok := fromRaw.(string)
	if !ok {
		return ToolResult{"rename_file", nil, "el argumento 'from' debe ser string"}
	}

	to, ok := toRaw.(string)
	if !ok {
		return ToolResult{"rename_file", nil, "el argumento 'to' debe ser string"}
	}

	fullFrom := filepath.Join("workspace", from)
	fullTo := filepath.Join("workspace", to)

	// Verificar existencia del origen
	if _, err := os.Stat(fullFrom); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"rename_file", nil, fmt.Sprintf("el archivo o directorio '%s' no existe", from)}
		}
		return ToolResult{"rename_file", nil, fmt.Sprintf("error accediendo al origen: %v", err)}
	}

	// Crear directorios destino si no existen
	if err := os.MkdirAll(filepath.Dir(fullTo), 0755); err != nil {
		return ToolResult{"rename_file", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Renombrar
	if err := os.Rename(fullFrom, fullTo); err != nil {
		return ToolResult{"rename_file", nil, fmt.Sprintf("error renombrando: %v", err)}
	}

	return ToolResult{
		ToolName: "rename_file",
		Result:   fmt.Sprintf("'%s' renombrado a '%s' correctamente", from, to),
	}
}

// copy_file: copia un archivo o directorio dentro del workspace
func copyFileTool(args map[string]interface{}) ToolResult {
	fromRaw, ok := args["from"]
	if !ok {
		return ToolResult{"copy_file", nil, "falta argumento obligatorio: from"}
	}

	toRaw, ok := args["to"]
	if !ok {
		return ToolResult{"copy_file", nil, "falta argumento obligatorio: to"}
	}

	from, ok := fromRaw.(string)
	if !ok {
		return ToolResult{"copy_file", nil, "el argumento 'from' debe ser string"}
	}

	to, ok := toRaw.(string)
	if !ok {
		return ToolResult{"copy_file", nil, "el argumento 'to' debe ser string"}
	}

	fullFrom := filepath.Join("workspace", from)
	fullTo := filepath.Join("workspace", to)

	// Verificar existencia del origen
	info, err := os.Stat(fullFrom)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"copy_file", nil, fmt.Sprintf("el archivo o directorio '%s' no existe", from)}
		}
		return ToolResult{"copy_file", nil, fmt.Sprintf("error accediendo al origen: %v", err)}
	}

	// Crear directorios destino si no existen
	if err := os.MkdirAll(filepath.Dir(fullTo), 0755); err != nil {
		return ToolResult{"copy_file", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Si es directorio → copiar recursivamente
	if info.IsDir() {
		err := filepath.Walk(fullFrom, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			rel, _ := filepath.Rel(fullFrom, path)
			dest := filepath.Join(fullTo, rel)

			if info.IsDir() {
				return os.MkdirAll(dest, 0755)
			}

			// Copiar archivo
			return copySingleFile(path, dest)
		})

		if err != nil {
			return ToolResult{"copy_file", nil, fmt.Sprintf("error copiando directorio: %v", err)}
		}

		return ToolResult{"copy_file", fmt.Sprintf("directorio '%s' copiado a '%s' correctamente", from, to), ""}
	}

	// Si es archivo → copiar archivo único
	if err := copySingleFile(fullFrom, fullTo); err != nil {
		return ToolResult{"copy_file", nil, fmt.Sprintf("error copiando archivo: %v", err)}
	}

	return ToolResult{"copy_file", fmt.Sprintf("archivo '%s' copiado a '%s' correctamente", from, to), ""}
}

// Función auxiliar para copiar un archivo individual
func copySingleFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

// move_file: mueve un archivo o directorio dentro del workspace
func moveFileTool(args map[string]interface{}) ToolResult {
	fromRaw, ok := args["from"]
	if !ok {
		return ToolResult{"move_file", nil, "falta argumento obligatorio: from"}
	}

	toRaw, ok := args["to"]
	if !ok {
		return ToolResult{"move_file", nil, "falta argumento obligatorio: to"}
	}

	from, ok := fromRaw.(string)
	if !ok {
		return ToolResult{"move_file", nil, "el argumento 'from' debe ser string"}
	}

	to, ok := toRaw.(string)
	if !ok {
		return ToolResult{"move_file", nil, "el argumento 'to' debe ser string"}
	}

	fullFrom := filepath.Join("workspace", from)
	fullTo := filepath.Join("workspace", to)

	// Verificar existencia del origen
	if _, err := os.Stat(fullFrom); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"move_file", nil, fmt.Sprintf("el archivo o directorio '%s' no existe", from)}
		}
		return ToolResult{"move_file", nil, fmt.Sprintf("error accediendo al origen: %v", err)}
	}

	// Crear directorios destino si no existen
	if err := os.MkdirAll(filepath.Dir(fullTo), 0755); err != nil {
		return ToolResult{"move_file", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Mover (rename)
	if err := os.Rename(fullFrom, fullTo); err != nil {
		return ToolResult{"move_file", nil, fmt.Sprintf("error moviendo archivo o directorio: %v", err)}
	}

	return ToolResult{
		ToolName: "move_file",
		Result:   fmt.Sprintf("'%s' movido a '%s' correctamente", from, to),
	}
}

// create_file: crea un archivo dentro del workspace, con contenido opcional
func createFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"create_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"create_file", nil, "el argumento 'path' debe ser string"}
	}

	content := ""
	if c, ok := args["content"]; ok {
		if cStr, ok := c.(string); ok {
			content = cStr
		}
	}

	fullPath := filepath.Join("workspace", path)

	// Crear directorios si no existen
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return ToolResult{"create_file", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Crear archivo con contenido
	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return ToolResult{"create_file", nil, fmt.Sprintf("error creando archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "create_file",
		Result:   fmt.Sprintf("archivo '%s' creado correctamente", path),
	}
}

// file_exists: verifica si un archivo o directorio existe dentro del workspace
func fileExistsTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"file_exists", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"file_exists", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{
				ToolName: "file_exists",
				Result: map[string]interface{}{
					"exists": false,
					"path":   path,
				},
			}
		}
		return ToolResult{"file_exists", nil, fmt.Sprintf("error comprobando archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "file_exists",
		Result: map[string]interface{}{
			"exists": true,
			"path":   path,
		},
	}
}

// read_dir: lista el contenido de un directorio con metadatos
func readDirTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"read_dir", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"read_dir", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return ToolResult{"read_dir", nil, fmt.Sprintf("error leyendo directorio: %v", err)}
	}

	var results []map[string]interface{}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"name":        entry.Name(),
			"is_dir":      entry.IsDir(),
			"size":        info.Size(),
			"modified":    info.ModTime().String(),
			"permissions": info.Mode().String(),
		})
	}

	return ToolResult{
		ToolName: "read_dir",
		Result:   results,
	}
}

// append_file: añade contenido al final de un archivo dentro del workspace
func appendFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"append_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"append_file", nil, "el argumento 'path' debe ser string"}
	}

	contentRaw, ok := args["content"]
	if !ok {
		return ToolResult{"append_file", nil, "falta argumento obligatorio: content"}
	}

	content, ok := contentRaw.(string)
	if !ok {
		return ToolResult{"append_file", nil, "el argumento 'content' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"append_file", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"append_file", nil, fmt.Sprintf("error accediendo al archivo: %v", err)}
	}

	// Abrir en modo append
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return ToolResult{"append_file", nil, fmt.Sprintf("error abriendo archivo: %v", err)}
	}
	defer f.Close()

	// Añadir contenido
	if _, err := f.WriteString(content); err != nil {
		return ToolResult{"append_file", nil, fmt.Sprintf("error escribiendo contenido: %v", err)}
	}

	return ToolResult{
		ToolName: "append_file",
		Result:   fmt.Sprintf("contenido añadido correctamente a '%s'", path),
	}
}

// truncate_file: vacía completamente un archivo dentro del workspace
func truncateFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"truncate_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"truncate_file", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"truncate_file", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"truncate_file", nil, fmt.Sprintf("error accediendo al archivo: %v", err)}
	}

	// Truncar archivo (dejarlo vacío)
	err := os.WriteFile(fullPath, []byte(""), 0644)
	if err != nil {
		return ToolResult{"truncate_file", nil, fmt.Sprintf("error truncando archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "truncate_file",
		Result:   fmt.Sprintf("archivo '%s' truncado correctamente", path),
	}
}

// stat_file: devuelve metadatos de un archivo o directorio dentro del workspace
func statFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"stat_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"stat_file", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"stat_file", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"stat_file", nil, fmt.Sprintf("error obteniendo metadatos: %v", err)}
	}

	result := map[string]interface{}{
		"path":        path,
		"is_dir":      info.IsDir(),
		"size":        info.Size(),
		"modified":    info.ModTime().String(),
		"permissions": info.Mode().String(),
	}

	return ToolResult{
		ToolName: "stat_file",
		Result:   result,
	}
}

// touch_file: crea un archivo vacío si no existe o actualiza su timestamp
func touchFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"touch_file", nil, "falta argumento obligatorio: path"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"touch_file", nil, "el argumento 'path' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Si el archivo no existe → crearlo vacío
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		// Crear directorios si no existen
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return ToolResult{"touch_file", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
		}

		err := os.WriteFile(fullPath, []byte(""), 0644)
		if err != nil {
			return ToolResult{"touch_file", nil, fmt.Sprintf("error creando archivo: %v", err)}
		}

		return ToolResult{
			ToolName: "touch_file",
			Result:   fmt.Sprintf("archivo '%s' creado correctamente", path),
		}
	}

	// Si existe → actualizar timestamp
	now := time.Now()
	if err := os.Chtimes(fullPath, now, now); err != nil {
		return ToolResult{"touch_file", nil, fmt.Sprintf("error actualizando timestamp: %v", err)}
	}

	return ToolResult{
		ToolName: "touch_file",
		Result:   fmt.Sprintf("timestamp de '%s' actualizado correctamente", path),
	}
}

// search_replace: busca y reemplaza texto dentro de un archivo del workspace
func searchReplaceTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"search_replace", nil, "falta argumento obligatorio: path"}
	}

	searchRaw, ok := args["search"]
	if !ok {
		return ToolResult{"search_replace", nil, "falta argumento obligatorio: search"}
	}

	replaceRaw, ok := args["replace"]
	if !ok {
		return ToolResult{"search_replace", nil, "falta argumento obligatorio: replace"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"search_replace", nil, "el argumento 'path' debe ser string"}
	}

	search, ok := searchRaw.(string)
	if !ok {
		return ToolResult{"search_replace", nil, "el argumento 'search' debe ser string"}
	}

	replace, ok := replaceRaw.(string)
	if !ok {
		return ToolResult{"search_replace", nil, "el argumento 'replace' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo
	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"search_replace", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"search_replace", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)

	// Contar ocurrencias
	count := strings.Count(content, search)

	// Reemplazar
	newContent := strings.ReplaceAll(content, search, replace)

	// Guardar archivo
	err = os.WriteFile(fullPath, []byte(newContent), 0644)
	if err != nil {
		return ToolResult{"search_replace", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "search_replace",
		Result: map[string]interface{}{
			"path":         path,
			"replacements": count,
			"search":       search,
			"replace":      replace,
			"success":      true,
		},
	}
}

// apply_patch_auto: aplica un parche inteligente sin requerir contexto exacto
func applyPatchAutoTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "falta argumento obligatorio: path"}
	}

	patchRaw, ok := args["patch"]
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "falta argumento obligatorio: patch"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "el argumento 'path' debe ser string"}
	}

	patch, ok := patchRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_auto", nil, "el argumento 'patch' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo original
	originalBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	original := string(originalBytes)

	// Dividir el parche en líneas
	lines := strings.Split(patch, "\n")

	var result strings.Builder
	result.WriteString(original)

	// Aplicación automática:
	// - Si la línea empieza con "+" → añadir al final
	// - Si empieza con "-" → eliminar todas las ocurrencias
	// - Si empieza con "~" → reemplazo inteligente: "~buscar => reemplazo"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "+"):
			// Añadir al final
			result.WriteString("\n" + strings.TrimPrefix(line, "+"))

		case strings.HasPrefix(line, "-"):
			// Eliminar todas las ocurrencias
			target := strings.TrimPrefix(line, "-")
			resultStr := result.String()
			resultStr = strings.ReplaceAll(resultStr, target, "")
			result.Reset()
			result.WriteString(resultStr)

		case strings.HasPrefix(line, "~"):
			// Reemplazo inteligente "~buscar => reemplazo"
			body := strings.TrimPrefix(line, "~")
			parts := strings.SplitN(body, "=>", 2)
			if len(parts) == 2 {
				search := strings.TrimSpace(parts[0])
				replace := strings.TrimSpace(parts[1])
				resultStr := result.String()
				resultStr = strings.ReplaceAll(resultStr, search, replace)
				result.Reset()
				result.WriteString(resultStr)
			}
		}
	}

	// Guardar archivo modificado
	err = os.WriteFile(fullPath, []byte(result.String()), 0644)
	if err != nil {
		return ToolResult{"apply_patch_auto", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "apply_patch_auto",
		Result:   fmt.Sprintf("parche inteligente aplicado correctamente a '%s'", path),
	}
}

// diff_files: compara dos archivos dentro del workspace y devuelve un diff estilo unified
func diffFilesTool(args map[string]interface{}) ToolResult {
	fromRaw, ok := args["from"]
	if !ok {
		return ToolResult{"diff_files", nil, "falta argumento obligatorio: from"}
	}

	toRaw, ok := args["to"]
	if !ok {
		return ToolResult{"diff_files", nil, "falta argumento obligatorio: to"}
	}

	from, ok := fromRaw.(string)
	if !ok {
		return ToolResult{"diff_files", nil, "el argumento 'from' debe ser string"}
	}

	to, ok := toRaw.(string)
	if !ok {
		return ToolResult{"diff_files", nil, "el argumento 'to' debe ser string"}
	}

	fullFrom := filepath.Join("workspace", from)
	fullTo := filepath.Join("workspace", to)

	// Leer archivos
	aBytes, err := os.ReadFile(fullFrom)
	if err != nil {
		return ToolResult{"diff_files", nil, fmt.Sprintf("error leyendo '%s': %v", from, err)}
	}

	bBytes, err := os.ReadFile(fullTo)
	if err != nil {
		return ToolResult{"diff_files", nil, fmt.Sprintf("error leyendo '%s': %v", to, err)}
	}

	aLines := strings.Split(string(aBytes), "\n")
	bLines := strings.Split(string(bBytes), "\n")

	// Generar diff estilo unified
	var diff strings.Builder
	diff.WriteString(fmt.Sprintf("--- %s\n", from))
	diff.WriteString(fmt.Sprintf("+++ %s\n", to))

	max := len(aLines)
	if len(bLines) > max {
		max = len(bLines)
	}

	for i := 0; i < max; i++ {
		var aLine, bLine string

		if i < len(aLines) {
			aLine = aLines[i]
		}
		if i < len(bLines) {
			bLine = bLines[i]
		}

		if aLine == bLine {
			diff.WriteString(" " + aLine + "\n")
		} else {
			if aLine != "" {
				diff.WriteString("-" + aLine + "\n")
			}
			if bLine != "" {
				diff.WriteString("+" + bLine + "\n")
			}
		}
	}

	return ToolResult{
		ToolName: "diff_files",
		Result:   diff.String(),
	}
}

// zip_dir: comprime un directorio dentro del workspace en un archivo .zip
func zipDirTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"zip_dir", nil, "falta argumento obligatorio: path"}
	}

	outRaw, ok := args["output"]
	if !ok {
		return ToolResult{"zip_dir", nil, "falta argumento obligatorio: output"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"zip_dir", nil, "el argumento 'path' debe ser string"}
	}

	output, ok := outRaw.(string)
	if !ok {
		return ToolResult{"zip_dir", nil, "el argumento 'output' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)
	fullOutput := filepath.Join("workspace", output)

	// Verificar que el directorio existe
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"zip_dir", nil, fmt.Sprintf("el directorio '%s' no existe", path)}
		}
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error accediendo al directorio: %v", err)}
	}

	if !info.IsDir() {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("'%s' no es un directorio", path)}
	}

	// Crear directorios destino si no existen
	if err := os.MkdirAll(filepath.Dir(fullOutput), 0755); err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error creando directorios destino: %v", err)}
	}

	// Crear archivo zip
	outFile, err := os.Create(fullOutput)
	if err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error creando archivo zip: %v", err)}
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Recorrer el directorio y añadir archivos
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(fullPath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return ToolResult{"zip_dir", nil, fmt.Sprintf("error comprimiendo directorio: %v", err)}
	}

	return ToolResult{
		ToolName: "zip_dir",
		Result:   fmt.Sprintf("directorio '%s' comprimido correctamente en '%s'", path, output),
	}
}

// unzip: descomprime un archivo .zip dentro del workspace
func unzipTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"unzip", nil, "falta argumento obligatorio: path"}
	}

	destRaw, ok := args["dest"]
	if !ok {
		return ToolResult{"unzip", nil, "falta argumento obligatorio: dest"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"unzip", nil, "el argumento 'path' debe ser string"}
	}

	dest, ok := destRaw.(string)
	if !ok {
		return ToolResult{"unzip", nil, "el argumento 'dest' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)
	fullDest := filepath.Join("workspace", dest)

	// Verificar que el archivo existe
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"unzip", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"unzip", nil, fmt.Sprintf("error accediendo al archivo: %v", err)}
	}

	// Abrir archivo zip
	r, err := zip.OpenReader(fullPath)
	if err != nil {
		return ToolResult{"unzip", nil, fmt.Sprintf("error abriendo zip: %v", err)}
	}
	defer r.Close()

	// Crear directorio destino
	if err := os.MkdirAll(fullDest, 0755); err != nil {
		return ToolResult{"unzip", nil, fmt.Sprintf("error creando directorio destino: %v", err)}
	}

	// Extraer archivos
	for _, f := range r.File {
		fpath := filepath.Join(fullDest, f.Name)

		// Evitar path traversal
		if !strings.HasPrefix(fpath, filepath.Clean(fullDest)+string(os.PathSeparator)) {
			return ToolResult{"unzip", nil, "zip contiene rutas inválidas (path traversal)"}
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}

		// Crear directorio si no existe
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error creando directorio: %v", err)}
		}

		// Extraer archivo
		rc, err := f.Open()
		if err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error leyendo archivo del zip: %v", err)}
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return ToolResult{"unzip", nil, fmt.Sprintf("error creando archivo: %v", err)}
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return ToolResult{"unzip", nil, fmt.Sprintf("error extrayendo archivo: %v", err)}
		}
	}

	return ToolResult{
		ToolName: "unzip",
		Result:   fmt.Sprintf("archivo '%s' descomprimido correctamente en '%s'", path, dest),
	}
}

// search_regex: busca coincidencias usando expresiones regulares dentro de un archivo
func searchRegexTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"search_regex", nil, "falta argumento obligatorio: path"}
	}

	regexRaw, ok := args["regex"]
	if !ok {
		return ToolResult{"search_regex", nil, "falta argumento obligatorio: regex"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"search_regex", nil, "el argumento 'path' debe ser string"}
	}

	regexStr, ok := regexRaw.(string)
	if !ok {
		return ToolResult{"search_regex", nil, "el argumento 'regex' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Leer archivo
	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"search_regex", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"search_regex", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)

	// Compilar regex
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return ToolResult{"search_regex", nil, fmt.Sprintf("regex inválida: %v", err)}
	}

	// Buscar coincidencias
	matches := re.FindAllStringIndex(content, -1)

	var results []map[string]interface{}
	for _, m := range matches {
		start := m[0]
		end := m[1]
		results = append(results, map[string]interface{}{
			"start": start,
			"end":   end,
			"match": content[start:end],
		})
	}

	return ToolResult{
		ToolName: "search_regex",
		Result: map[string]interface{}{
			"path":    path,
			"regex":   regexStr,
			"matches": results,
			"count":   len(results),
		},
	}
}

// apply_patch_structured: modificaciones semánticas de alto nivel
func applyPatchStructuredTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "falta argumento obligatorio: path"}
	}

	opRaw, ok := args["op"]
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "falta argumento obligatorio: op"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "el argumento 'path' debe ser string"}
	}

	op, ok := opRaw.(string)
	if !ok {
		return ToolResult{"apply_patch_structured", nil, "el argumento 'op' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("el archivo '%s' no existe", path)}
		}
		return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("error leyendo archivo: %v", err)}
	}

	content := string(contentBytes)
	updated := content

	switch op {

	// ---------------------------------------------------------
	// Insertar import
	// ---------------------------------------------------------
	case "insert_import":
		importRaw := args["import"]
		if importRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "falta argumento 'import'"}
		}
		importLine, _ := importRaw.(string)

		re := regexp.MustCompile(`(?m)^import\s*\(`)
		if re.MatchString(updated) {
			updated = re.ReplaceAllString(updated, "import (\n    "+importLine)
		} else {
			updated = "import (\n    " + importLine + "\n)\n\n" + updated
		}

	// ---------------------------------------------------------
	// Insertar antes de una función
	// ---------------------------------------------------------
	case "insert_before_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?m)^func\s+` + regexp.QuoteMeta(name) + `\s*\(`)
		loc := re.FindStringIndex(updated)
		if loc == nil {
			return ToolResult{"apply_patch_structured", nil, "función no encontrada"}
		}

		updated = updated[:loc[0]] + code + "\n" + updated[loc[0]:]

	// ---------------------------------------------------------
	// Insertar después de una función
	// ---------------------------------------------------------
	case "insert_after_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		match := re.FindStringIndex(updated)
		if match == nil {
			return ToolResult{"apply_patch_structured", nil, "función no encontrada"}
		}

		updated = updated[:match[1]] + "\n" + code + "\n" + updated[match[1]:]

	// ---------------------------------------------------------
	// Reemplazar función completa
	// ---------------------------------------------------------
	case "replace_func":
		nameRaw := args["name"]
		codeRaw := args["code"]
		if nameRaw == nil || codeRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'name' y/o 'code'"}
		}
		name := nameRaw.(string)
		code := codeRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		updated = re.ReplaceAllString(updated, code)

	// ---------------------------------------------------------
	// Eliminar función completa
	// ---------------------------------------------------------
	case "delete_func":
		nameRaw := args["name"]
		if nameRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "falta argumento 'name'"}
		}
		name := nameRaw.(string)

		re := regexp.MustCompile(`(?s)func\s+` + regexp.QuoteMeta(name) + `\s*\([^)]*\)\s*{.*?}`)
		updated = re.ReplaceAllString(updated, "")

	// ---------------------------------------------------------
	// Reemplazo por regex
	// ---------------------------------------------------------
	case "regex_replace":
		regexRaw := args["regex"]
		replaceRaw := args["replace"]
		if regexRaw == nil || replaceRaw == nil {
			return ToolResult{"apply_patch_structured", nil, "faltan argumentos 'regex' y/o 'replace'"}
		}
		regex := regexRaw.(string)
		replace := replaceRaw.(string)

		re, err := regexp.Compile(regex)
		if err != nil {
			return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("regex inválida: %v", err)}
		}

		updated = re.ReplaceAllString(updated, replace)

	default:
		return ToolResult{"apply_patch_structured", nil, "operación desconocida"}
	}

	// Guardar archivo
	err = os.WriteFile(fullPath, []byte(updated), 0644)
	if err != nil {
		return ToolResult{"apply_patch_structured", nil, fmt.Sprintf("error escribiendo archivo: %v", err)}
	}

	return ToolResult{
		ToolName: "apply_patch_structured",
		Result:   fmt.Sprintf("parche estructurado aplicado correctamente a '%s'", path),
	}
}

package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// readFileTool
// writeFileTool
// createFileTool
// deleteFileTool
// renameFileTool
// copyFileTool
// moveFileTool
// fileExistsTool
// readDirTool
// statFileTool
// touchFileTool

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

package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

// search_regex_multi: busca coincidencias regex en múltiples archivos dentro de un directorio
func searchRegexMultiTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"search_regex_multi", nil, "falta argumento obligatorio: path"}
	}

	regexRaw, ok := args["regex"]
	if !ok {
		return ToolResult{"search_regex_multi", nil, "falta argumento obligatorio: regex"}
	}

	path, ok := pathRaw.(string)
	if !ok {
		return ToolResult{"search_regex_multi", nil, "el argumento 'path' debe ser string"}
	}

	regexStr, ok := regexRaw.(string)
	if !ok {
		return ToolResult{"search_regex_multi", nil, "el argumento 'regex' debe ser string"}
	}

	fullPath := filepath.Join("workspace", path)

	// Verificar que el directorio existe
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ToolResult{"search_regex_multi", nil, fmt.Sprintf("el directorio '%s' no existe", path)}
		}
		return ToolResult{"search_regex_multi", nil, fmt.Sprintf("error accediendo al directorio: %v", err)}
	}

	if !info.IsDir() {
		return ToolResult{"search_regex_multi", nil, fmt.Sprintf("'%s' no es un directorio", path)}
	}

	// Compilar regex
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return ToolResult{"search_regex_multi", nil, fmt.Sprintf("regex inválida: %v", err)}
	}

	results := map[string][]map[string]interface{}{}

	// Recorrer recursivamente
	err = filepath.Walk(fullPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar directorios
		if info.IsDir() {
			return nil
		}

		// Leer archivo
		contentBytes, err := os.ReadFile(p)
		if err != nil {
			return nil // ignorar errores por archivo individual
		}

		content := string(contentBytes)
		matches := re.FindAllStringIndex(content, -1)

		if len(matches) > 0 {
			rel, _ := filepath.Rel(fullPath, p)

			for _, m := range matches {
				start := m[0]
				end := m[1]

				results[rel] = append(results[rel], map[string]interface{}{
					"start": start,
					"end":   end,
					"match": content[start:end],
				})
			}
		}

		return nil
	})

	if err != nil {
		return ToolResult{"search_regex_multi", nil, fmt.Sprintf("error recorriendo directorio: %v", err)}
	}

	return ToolResult{
		ToolName: "search_regex_multi",
		Result: map[string]interface{}{
			"path":    path,
			"regex":   regexStr,
			"results": results,
			"files":   len(results),
		},
	}
}

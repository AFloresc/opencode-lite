package tools

import (
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// extractFunctionsTool
// Extrae funciones de un archivo Go (nombre, firma, cuerpo)
// ------------------------------------------------------------
func extractFunctionsTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"extract_functions", nil, "falta argumento obligatorio: path"}
	}

	path, _ := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"extract_functions", nil, err.Error()}
	}

	// Regex para capturar funciones completas
	re := regexp.MustCompile(`(?ms)^func\s+([A-Za-z0-9_]+)\s*\((.*?)\)\s*(.*?)\{(.*?)\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	var funcs []map[string]interface{}

	for _, m := range matches {
		funcs = append(funcs, map[string]interface{}{
			"name":   m[1],
			"params": m[2],
			"header": m[3],
			"body":   strings.TrimSpace(m[4]),
		})
	}

	return ToolResult{"extract_functions", map[string]interface{}{
		"path":      path,
		"functions": funcs,
		"count":     len(funcs),
	}, ""}
}

// ------------------------------------------------------------
// extractTypesTool
// Extrae structs, interfaces, alias y tipos definidos
// ------------------------------------------------------------
func extractTypesTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"extract_types", nil, "falta argumento obligatorio: path"}
	}

	path, _ := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"extract_types", nil, err.Error()}
	}

	re := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+(struct|interface|=)`)
	matches := re.FindAllStringSubmatch(content, -1)

	var types []map[string]interface{}

	for _, m := range matches {
		types = append(types, map[string]interface{}{
			"name": m[1],
			"kind": m[2],
		})
	}

	return ToolResult{"extract_types", map[string]interface{}{
		"path":  path,
		"types": types,
		"count": len(types),
	}, ""}
}

// ------------------------------------------------------------
// extractCommentsBlockTool
// Extrae comentarios multilinea /* ... */
// ------------------------------------------------------------
func extractCommentsBlockTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"extract_comments_block", nil, "falta argumento obligatorio: path"}
	}

	path, _ := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"extract_comments_block", nil, err.Error()}
	}

	re := regexp.MustCompile(`(?s)/\*.*?\*/`)
	matches := re.FindAllString(content, -1)

	return ToolResult{"extract_comments_block", map[string]interface{}{
		"path":     path,
		"comments": matches,
		"count":    len(matches),
	}, ""}
}

// ------------------------------------------------------------
// semanticIndexTool
// Crea un índice semántico del archivo:
// - funciones
// - structs
// - interfaces
// - alias
// - comentarios
// ------------------------------------------------------------
func semanticIndexTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"semantic_index", nil, "falta argumento obligatorio: path"}
	}

	path, _ := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"semantic_index", nil, err.Error()}
	}

	// Funciones
	reFunc := regexp.MustCompile(`(?m)^func\s+([A-Za-z0-9_]+)\s*\(`)
	funcs := reFunc.FindAllStringSubmatch(content, -1)
	var funcNames []string
	for _, f := range funcs {
		funcNames = append(funcNames, f[1])
	}

	// Structs
	reStruct := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+struct`)
	structs := reStruct.FindAllStringSubmatch(content, -1)
	var structNames []string
	for _, s := range structs {
		structNames = append(structNames, s[1])
	}

	// Interfaces
	reIface := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+interface`)
	ifaces := reIface.FindAllStringSubmatch(content, -1)
	var ifaceNames []string
	for _, i := range ifaces {
		ifaceNames = append(ifaceNames, i[1])
	}

	// Alias
	reAlias := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+=`)
	aliases := reAlias.FindAllStringSubmatch(content, -1)
	var aliasNames []string
	for _, a := range aliases {
		aliasNames = append(aliasNames, a[1])
	}

	// Comentarios multilinea
	reBlock := regexp.MustCompile(`(?s)/\*.*?\*/`)
	comments := reBlock.FindAllString(content, -1)

	return ToolResult{"semantic_index", map[string]interface{}{
		"path":       path,
		"functions":  funcNames,
		"structs":    structNames,
		"interfaces": ifaceNames,
		"aliases":    aliasNames,
		"comments":   comments,
	}, ""}
}

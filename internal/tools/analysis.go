package tools

import (
	"path/filepath"
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// analysisDependenciesTool
// Analiza imports y dependencias entre archivos Go
// ------------------------------------------------------------
func analysisDependenciesTool(args map[string]interface{}) ToolResult {
	rootRaw, ok := args["root"]
	if !ok {
		return ToolResult{"analysis_dependencies", nil, "falta argumento obligatorio: root"}
	}

	root := rootRaw.(string)

	files, err := listFilesRecursive(root)
	if err != nil {
		return ToolResult{"analysis_dependencies", nil, err.Error()}
	}

	depMap := map[string][]string{}

	reImport := regexp.MustCompile(`import\s+"([^"]+)"`)
	reImportBlock := regexp.MustCompile(`import\s*\((?s)(.*?)\)`)

	for _, f := range files {
		if !strings.HasSuffix(f, ".go") {
			continue
		}

		full := filepath.Join(root, f)
		content, err := readFile(full)
		if err != nil {
			continue
		}

		var imports []string

		// import "x"
		for _, m := range reImport.FindAllStringSubmatch(content, -1) {
			imports = append(imports, m[1])
		}

		// import ( "x" "y" )
		for _, block := range reImportBlock.FindAllStringSubmatch(content, -1) {
			lines := strings.Split(block[1], "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if strings.HasPrefix(l, "\"") && strings.HasSuffix(l, "\"") {
					imports = append(imports, strings.Trim(l, "\""))
				}
			}
		}

		depMap[f] = imports
	}

	return ToolResult{"analysis_dependencies", map[string]interface{}{
		"root":         root,
		"dependencies": depMap,
	}, ""}
}

// ------------------------------------------------------------
// analysisCyclomaticTool
// Calcula complejidad ciclomática básica por función
// ------------------------------------------------------------
func analysisCyclomaticTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"analysis_cyclomatic", nil, "falta argumento obligatorio: path"}
	}

	path := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"analysis_cyclomatic", nil, err.Error()}
	}

	// Detectar funciones
	reFunc := regexp.MustCompile(`(?m)^func\s+([A-Za-z0-9_]+)\s*\(`)
	funcs := reFunc.FindAllStringSubmatchIndex(content, -1)

	results := []map[string]interface{}{}

	for i, f := range funcs {
		name := content[f[2]:f[3]]

		// Determinar rango de la función
		start := f[0]
		end := len(content)
		if i < len(funcs)-1 {
			end = funcs[i+1][0]
		}

		body := content[start:end]

		// Contar decisiones
		score := 1
		keywords := []string{"if ", "for ", "switch ", "case ", "&&", "||"}

		for _, k := range keywords {
			score += strings.Count(body, k)
		}

		results = append(results, map[string]interface{}{
			"function":   name,
			"complexity": score,
		})
	}

	return ToolResult{"analysis_cyclomatic", map[string]interface{}{
		"path":    path,
		"results": results,
	}, ""}
}

// ------------------------------------------------------------
// analysisDeadCodeTool
// Heurística simple para detectar funciones no usadas
// ------------------------------------------------------------
func analysisDeadCodeTool(args map[string]interface{}) ToolResult {
	rootRaw, ok := args["root"]
	if !ok {
		return ToolResult{"analysis_dead_code", nil, "falta argumento obligatorio: root"}
	}

	root := rootRaw.(string)

	files, err := listFilesRecursive(root)
	if err != nil {
		return ToolResult{"analysis_dead_code", nil, err.Error()}
	}

	// 1. Extraer todas las funciones
	allFuncs := map[string]string{} // name → file

	reFunc := regexp.MustCompile(`(?m)^func\s+([A-Za-z0-9_]+)\s*\(`)

	for _, f := range files {
		if !strings.HasSuffix(f, ".go") {
			continue
		}

		full := filepath.Join(root, f)
		content, err := readFile(full)
		if err != nil {
			continue
		}

		for _, m := range reFunc.FindAllStringSubmatch(content, -1) {
			allFuncs[m[1]] = f
		}
	}

	// 2. Buscar referencias
	used := map[string]bool{}

	for _, f := range files {
		full := filepath.Join(root, f)
		content, err := readFile(full)
		if err != nil {
			continue
		}

		for name := range allFuncs {
			if strings.Contains(content, name+"(") {
				used[name] = true
			}
		}
	}

	// 3. Dead code = funciones nunca llamadas
	dead := []map[string]interface{}{}

	for name, file := range allFuncs {
		if !used[name] {
			dead = append(dead, map[string]interface{}{
				"function": name,
				"file":     file,
			})
		}
	}

	return ToolResult{"analysis_dead_code", map[string]interface{}{
		"root":  root,
		"dead":  dead,
		"count": len(dead),
	}, ""}
}

// ------------------------------------------------------------
// analysisMetricsTool
// Métricas globales del archivo:
// - LOC
// - líneas de código reales
// - comentarios
// - funciones
// - structs
// - interfaces
// ------------------------------------------------------------
func analysisMetricsTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"analysis_metrics", nil, "falta argumento obligatorio: path"}
	}

	path := pathRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"analysis_metrics", nil, err.Error()}
	}

	lines := strings.Split(content, "\n")

	loc := len(lines)
	codeLines := 0
	commentLines := 0

	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "//") {
			commentLines++
		} else {
			codeLines++
		}
	}

	// Funciones
	reFunc := regexp.MustCompile(`(?m)^func\s+([A-Za-z0-9_]+)\s*\(`)
	funcs := reFunc.FindAllStringSubmatch(content, -1)

	// Structs
	reStruct := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+struct`)
	structs := reStruct.FindAllStringSubmatch(content, -1)

	// Interfaces
	reIface := regexp.MustCompile(`(?m)^type\s+([A-Za-z0-9_]+)\s+interface`)
	ifaces := reIface.FindAllStringSubmatch(content, -1)

	return ToolResult{"analysis_metrics", map[string]interface{}{
		"path":          path,
		"loc":           loc,
		"code_lines":    codeLines,
		"comment_lines": commentLines,
		"functions":     len(funcs),
		"structs":       len(structs),
		"interfaces":    len(ifaces),
	}, ""}
}

package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// refactorRenameSymbolTool
// Renombra un símbolo (función, variable, struct, interface)
// en TODO el proyecto (búsqueda exacta + heurística básica)
// ------------------------------------------------------------
func refactorRenameSymbolTool(args map[string]interface{}) ToolResult {
	oldRaw, ok := args["old"]
	if !ok {
		return ToolResult{"refactor_rename_symbol", nil, "falta argumento obligatorio: old"}
	}

	newRaw, ok := args["new"]
	if !ok {
		return ToolResult{"refactor_rename_symbol", nil, "falta argumento obligatorio: new"}
	}

	rootRaw, ok := args["root"]
	if !ok {
		return ToolResult{"refactor_rename_symbol", nil, "falta argumento obligatorio: root"}
	}

	oldName := oldRaw.(string)
	newName := newRaw.(string)
	root := rootRaw.(string)

	files, err := listFilesRecursive(root)
	if err != nil {
		return ToolResult{"refactor_rename_symbol", nil, err.Error()}
	}

	changed := []string{}

	// Regex para coincidencia exacta de símbolo
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(oldName) + `\b`)

	for _, f := range files {
		full := filepath.Join(root, f)
		content, err := readFile(full)
		if err != nil {
			continue
		}

		newContent := re.ReplaceAllString(content, newName)
		if newContent != content {
			writeFile(full, newContent)
			changed = append(changed, f)
		}
	}

	return ToolResult{"refactor_rename_symbol", map[string]interface{}{
		"old":     oldName,
		"new":     newName,
		"changed": changed,
		"count":   len(changed),
	}, ""}
}

// ------------------------------------------------------------
// refactorMoveFileTool
// Mueve un archivo y actualiza imports en todo el proyecto
// ------------------------------------------------------------
func refactorMoveFileTool(args map[string]interface{}) ToolResult {
	fromRaw, ok := args["from"]
	if !ok {
		return ToolResult{"refactor_move_file", nil, "falta argumento obligatorio: from"}
	}

	toRaw, ok := args["to"]
	if !ok {
		return ToolResult{"refactor_move_file", nil, "falta argumento obligatorio: to"}
	}

	rootRaw, ok := args["root"]
	if !ok {
		return ToolResult{"refactor_move_file", nil, "falta argumento obligatorio: root"}
	}

	from := fromRaw.(string)
	to := toRaw.(string)
	root := rootRaw.(string)

	fullFrom, err := safeJoinWorkspace(from)
	if err != nil {
		return ToolResult{"refactor_move_file", nil, err.Error()}
	}

	fullTo, err := safeJoinWorkspace(to)
	if err != nil {
		return ToolResult{"refactor_move_file", nil, err.Error()}
	}

	// Crear directorio destino si no existe
	os.MkdirAll(filepath.Dir(fullTo), 0755)

	// Mover archivo
	if err := os.Rename(fullFrom, fullTo); err != nil {
		return ToolResult{"refactor_move_file", nil, err.Error()}
	}

	// Actualizar imports en todo el proyecto
	files, err := listFilesRecursive(root)
	if err != nil {
		return ToolResult{"refactor_move_file", nil, err.Error()}
	}

	oldImport := filepath.Dir(from)
	newImport := filepath.Dir(to)

	changed := []string{}

	for _, f := range files {
		full := filepath.Join(root, f)
		content, err := readFile(full)
		if err != nil {
			continue
		}

		newContent := strings.ReplaceAll(content, oldImport, newImport)
		if newContent != content {
			writeFile(full, newContent)
			changed = append(changed, f)
		}
	}

	return ToolResult{"refactor_move_file", map[string]interface{}{
		"from":    from,
		"to":      to,
		"updated": changed,
		"count":   len(changed),
	}, ""}
}

// ------------------------------------------------------------
// refactorSplitFileTool
// Divide un archivo en varios según delimitadores o patrones
// ------------------------------------------------------------
func refactorSplitFileTool(args map[string]interface{}) ToolResult {
	pathRaw, ok := args["path"]
	if !ok {
		return ToolResult{"refactor_split_file", nil, "falta argumento obligatorio: path"}
	}

	patternRaw, ok := args["pattern"]
	if !ok {
		return ToolResult{"refactor_split_file", nil, "falta argumento obligatorio: pattern"}
	}

	path := pathRaw.(string)
	pattern := patternRaw.(string)

	content, err := readFile(path)
	if err != nil {
		return ToolResult{"refactor_split_file", nil, err.Error()}
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return ToolResult{"refactor_split_file", nil, "regex inválida: " + err.Error()}
	}

	parts := re.Split(content, -1)
	if len(parts) < 2 {
		return ToolResult{"refactor_split_file", nil, "no se encontraron divisiones"}
	}

	base := strings.TrimSuffix(path, filepath.Ext(path))
	ext := filepath.Ext(path)

	var outFiles []string

	for i, p := range parts {
		newPath := fmt.Sprintf("%s_part%d%s", base, i+1, ext)
		if err := writeFile(newPath, strings.TrimSpace(p)); err != nil {
			return ToolResult{"refactor_split_file", nil, err.Error()}
		}
		outFiles = append(outFiles, newPath)
	}

	return ToolResult{"refactor_split_file", map[string]interface{}{
		"path":  path,
		"parts": outFiles,
		"count": len(outFiles),
	}, ""}
}

// ------------------------------------------------------------
// refactorMergeFilesTool
// Fusiona varios archivos en uno solo
// ------------------------------------------------------------
func refactorMergeFilesTool(args map[string]interface{}) ToolResult {
	filesRaw, ok := args["files"]
	if !ok {
		return ToolResult{"refactor_merge_files", nil, "falta argumento obligatorio: files"}
	}

	outRaw, ok := args["out"]
	if !ok {
		return ToolResult{"refactor_merge_files", nil, "falta argumento obligatorio: out"}
	}

	files := filesRaw.([]interface{})
	outPath := outRaw.(string)

	var builder []string

	for _, f := range files {
		path := f.(string)
		content, err := readFile(path)
		if err != nil {
			return ToolResult{"refactor_merge_files", nil, err.Error()}
		}

		builder = append(builder,
			"// ---- BEGIN "+path+" ----",
			content,
			"// ---- END "+path+" ----",
			"",
		)
	}

	merged := strings.Join(builder, "\n")

	if err := writeFile(outPath, merged); err != nil {
		return ToolResult{"refactor_merge_files", nil, err.Error()}
	}

	return ToolResult{"refactor_merge_files", map[string]interface{}{
		"out":   outPath,
		"files": files,
		"ok":    true,
	}, ""}
}

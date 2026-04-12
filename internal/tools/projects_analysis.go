package tools

import (
	"os"
	"path/filepath"
	"sort"
)

// ------------------------------------------------------------
// projectStatsTool
// Devuelve número total de archivos y directorios en workspace/
// ------------------------------------------------------------
func projectStatsTool(args map[string]interface{}) ToolResult {
	root := "workspace"

	var files, dirs int

	filepath.Walk(root, func(_ string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			dirs++
		} else {
			files++
		}
		return nil
	})

	return ToolResult{"project_stats", map[string]interface{}{
		"files": files,
		"dirs":  dirs,
	}, ""}
}

// ------------------------------------------------------------
// largestFilesTool
// Devuelve los 10 archivos más grandes del workspace
// ------------------------------------------------------------
func largestFilesTool(args map[string]interface{}) ToolResult {
	root := "workspace"

	type entry struct {
		Path string
		Size int64
	}

	var list []entry

	filepath.Walk(root, func(p string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			list = append(list, entry{p, info.Size()})
		}
		return nil
	})

	sort.Slice(list, func(a, b int) bool { return list[a].Size > list[b].Size })

	top := list
	if len(top) > 10 {
		top = top[:10]
	}

	var out []map[string]interface{}
	for _, e := range top {
		rel, _ := filepath.Rel(root, e.Path)
		out = append(out, map[string]interface{}{
			"path": rel,
			"size": e.Size,
		})
	}

	return ToolResult{"largest_files", map[string]interface{}{
		"largest": out,
	}, ""}
}

// ------------------------------------------------------------
// fileTreeTool
// Devuelve un árbol completo del workspace (todos los paths)
// ------------------------------------------------------------
func fileTreeTool(args map[string]interface{}) ToolResult {
	root := "workspace"

	var tree []string

	filepath.Walk(root, func(p string, info os.FileInfo, _ error) error {
		rel, _ := filepath.Rel(root, p)
		tree = append(tree, rel)
		return nil
	})

	return ToolResult{"file_tree", map[string]interface{}{
		"tree": tree,
	}, ""}
}

package tools

import (
	"os"
	"path/filepath"
	"sort"
)

func runCmdProjectStats(main string, argsList []string) ToolResult {
	root := filepath.Join("workspace")
	var files, dirs int
	filepath.Walk(root, func(_ string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			dirs++
		} else {
			files++
		}
		return nil
	})
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"files":   files,
		"dirs":    dirs,
	}, ""}
}

func runCmdLargestFiles(main string, argsList []string) ToolResult {
	root := filepath.Join("workspace")
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
	out := []map[string]interface{}{}
	for _, e := range top {
		rel, _ := filepath.Rel(root, e.Path)
		out = append(out, map[string]interface{}{
			"path": rel,
			"size": e.Size,
		})
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"largest": out,
	}, ""}
}

func runCmdFileTree(main string, argsList []string) ToolResult {
	root := filepath.Join("workspace")
	var tree []string
	filepath.Walk(root, func(p string, info os.FileInfo, _ error) error {
		rel, _ := filepath.Rel(root, p)
		tree = append(tree, rel)
		return nil
	})
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"tree":    tree,
	}, ""}
}

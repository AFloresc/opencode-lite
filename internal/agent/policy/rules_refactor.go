package agent

func refactorRules() []Rule {
	return []Rule{
		{
			Match: func(goal string) bool {
				return containsAny(goal, "renombrar", "rename")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				symbol, newName := extractRename(goal)
				return "refactor_rename_symbol", map[string]interface{}{
					"symbol":   symbol,
					"new_name": newName,
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "mover archivo", "move file")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				file := extractFile(goal)
				newPath := extractNewPath(goal)
				return "refactor_move_file", map[string]interface{}{
					"path":     file,
					"new_path": newPath,
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "split", "dividir archivo")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				file := extractFile(goal)
				return "refactor_split_file", map[string]interface{}{
					"path": file,
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "merge", "fusionar archivos")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				files := extractFiles(goal)
				return "refactor_merge_files", map[string]interface{}{
					"files": files,
				}, false
			},
		},
	}
}

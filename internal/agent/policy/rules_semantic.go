package agent

func semanticRules() []Rule {
	return []Rule{
		{
			Match: func(goal string) bool {
				return containsAny(goal, "explica", "explain")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "summarize_file", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "resumir", "resume", "summary")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "summarize_file", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "detectar lenguaje", "detect language")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "detect_language", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "extraer funciones", "extract functions")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "extract_functions", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "extraer tipos", "extract types")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "extract_types", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "extraer comentarios", "extract comments")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "extract_comments_block", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "indexar", "semantic index")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "semantic_index", map[string]interface{}{
					"root": "workspace",
				}, false
			},
		},
	}
}

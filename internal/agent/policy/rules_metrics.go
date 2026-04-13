package agent

func metricsRules() []Rule {
	return []Rule{
		{
			Match: func(goal string) bool {
				return containsAny(goal, "métricas", "metrics", "estadísticas", "stats")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "analysis_metrics", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "dependencias", "dependencies")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "analysis_dependencies", map[string]interface{}{
					"root": "workspace",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "complejidad", "cyclomatic")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "analysis_cyclomatic", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "dead code", "código muerto")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "analysis_dead_code", map[string]interface{}{
					"root": "workspace",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "archivos grandes", "large files")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "largest_files", map[string]interface{}{
					"root": "workspace",
				}, false
			},
		},
	}
}

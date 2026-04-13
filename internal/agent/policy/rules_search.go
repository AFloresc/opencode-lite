package agent

func searchRules() []Rule {
	return []Rule{
		{
			Match: func(goal string) bool {
				return containsAny(goal, "buscar", "search", "grep", "encontrar")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": extractPattern(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "todos", "fixme", "pendientes")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "TODO|FIXME",
				}, false
			},
		},
	}
}

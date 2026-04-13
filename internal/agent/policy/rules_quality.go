package agent

func qualityRules() []Rule {
	return []Rule{
		// Formato / limpieza
		{
			Match: func(goal string) bool {
				return containsAny(goal, "formatear", "format code", "formato")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "format_code", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "limpiar imports", "clean imports")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "format_code", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},

		// Lint / sintaxis
		{
			Match: func(goal string) bool {
				return containsAny(goal, "validar sintaxis", "syntax check", "lint")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "lint_code", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "variables no usadas", "unused vars")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "lint_code", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "imports no usados", "unused imports")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "lint_code", map[string]interface{}{
					"path": extractFile(goal),
				}, false
			},
		},

		// Calidad estructural
		{
			Match: func(goal string) bool {
				return containsAny(goal, "funciones largas", "long functions")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "func .*\\{[\\s\\S]{200,}\\}",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "demasiados parámetros", "too many parameters")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "func [A-Za-z0-9_]+\\([^)]{40,}\\)",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "nombres malos", "bad names", "nombres poco descriptivos")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "\\b[a-zA-Z]{1,2}\\b",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "duplicación", "duplicate code", "duplicado")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "func .*\\{[\\s\\S]{100,}\\}",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "funciones sin comentarios", "undocumented functions")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "func [A-Za-z0-9_]+\\([^)]*\\) \\{",
				}, false
			},
		},
		{
			Match: func(goal string) bool {
				return containsAny(goal, "sin tests", "missing tests")
			},
			Apply: func(goal string) (string, map[string]interface{}, bool) {
				return "search_regex_multi", map[string]interface{}{
					"path":    "workspace",
					"pattern": "_test\\.go",
				}, false
			},
		},
	}
}

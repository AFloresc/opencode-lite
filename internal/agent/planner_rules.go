package agent

import "strings"

type QuickRule struct {
	Match func(goal string) bool
	Plan  func(goal string) Plan
}

func quickRules() []QuickRule {
	return []QuickRule{
		// === Análisis de proyecto ===
		{
			Match: func(g string) bool {
				return containsAny(g, "analizar proyecto", "project analysis", "auditar proyecto")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"listar archivos"},
						{"calcular métricas"},
						{"detectar dependencias"},
						{"buscar duplicación"},
						{"buscar funciones largas"},
					},
				}
			},
		},

		// === Explicar archivo ===
		{
			Match: func(g string) bool {
				return containsAny(g, "explicar archivo", "explain file", "entender archivo")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"extraer funciones"},
						{"extraer tipos"},
						{"extraer comentarios"},
						{"resumir archivo"},
					},
				}
			},
		},

		// === Limpieza ===
		{
			Match: func(g string) bool {
				return containsAny(g, "limpiar proyecto", "clean project", "formatear todo")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"listar archivos"},
						{"limpiar imports"},
						{"formatear"},
						{"buscar código muerto"},
					},
				}
			},
		},

		// === Refactor ===
		{
			Match: func(g string) bool {
				return containsAny(g, "refactor", "optimizar", "mejorar código")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"buscar funciones largas"},
						{"buscar duplicación"},
						{"buscar nombres poco descriptivos"},
						{"buscar demasiados parámetros"},
					},
				}
			},
		},

		// === Seguridad ===
		{
			Match: func(g string) bool {
				return containsAny(g, "seguridad", "security", "vulnerabilidad")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"buscar uso de exec.Command"},
						{"buscar rutas hardcodeadas"},
						{"buscar uso de md5 o sha1"},
						{"buscar errores ignorados"},
					},
				}
			},
		},

		// === Arquitectura ===
		{
			Match: func(g string) bool {
				return containsAny(g, "arquitectura", "architecture", "dependencias")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"detectar dependencias"},
						{"buscar ciclos de dependencias"},
						{"buscar paquetes demasiado grandes"},
					},
				}
			},
		},

		// === Documentación ===
		{
			Match: func(g string) bool {
				return containsAny(g, "documentar", "docs", "generar documentación")
			},
			Plan: func(g string) Plan {
				return Plan{
					Steps: []PlanStep{
						{"extraer comentarios"},
						{"extraer funciones"},
						{"generar documentación"},
					},
				}
			},
		},
	}
}

func matchQuickRule(goal string) (Plan, bool) {
	g := strings.ToLower(goal)
	for _, r := range quickRules() {
		if r.Match(g) {
			return r.Plan(g), true
		}
	}
	return Plan{}, false
}

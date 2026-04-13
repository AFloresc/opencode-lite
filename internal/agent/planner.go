package agent

import (
	"strings"
)

type PlanStep struct {
	Description string
}

type Plan struct {
	Steps []PlanStep
}

type Planner interface {
	MakePlan(goal string) Plan
}

type AdvancedHeuristicPlanner struct{}

func NewAdvancedHeuristicPlanner() *AdvancedHeuristicPlanner {
	return &AdvancedHeuristicPlanner{}
}

func (p *AdvancedHeuristicPlanner) MakePlan(goal string) Plan {
	g := strings.ToLower(goal)

	// Clasificación del goal
	switch {
	case containsAny(g, "analizar proyecto", "project analysis", "auditar proyecto"):
		return p.planProjectAnalysis()

	case containsAny(g, "explicar archivo", "explain file"):
		return p.planExplainFile()

	case containsAny(g, "limpiar proyecto", "clean project"):
		return p.planCleanup()

	case containsAny(g, "refactor", "optimizar", "mejorar código"):
		return p.planRefactor()

	case containsAny(g, "seguridad", "security"):
		return p.planSecurityAudit()

	case containsAny(g, "arquitectura", "architecture"):
		return p.planArchitecture()

	case containsAny(g, "documentar", "docs"):
		return p.planDocumentation()
	}

	// Plan genérico
	return Plan{
		Steps: []PlanStep{
			{Description: goal},
		},
	}
}

//
// ============================
// SUBPLANES CON HEURÍSTICAS
// ============================
//

func (p *AdvancedHeuristicPlanner) planProjectAnalysis() Plan {
	// Heurística: analizar tamaño del proyecto
	stats := AnalyzeProjectSize()

	steps := []PlanStep{
		{"listar archivos"},
		{"calcular métricas"},
		{"detectar dependencias"},
	}

	// Si hay muchos archivos, añadir pasos extra
	if stats.FileCount > 200 {
		steps = append(steps,
			PlanStep{"buscar duplicación"},
			PlanStep{"buscar funciones largas"},
			PlanStep{"buscar imports no usados"},
		)
	}

	// Si hay muchos paquetes, revisar arquitectura
	if stats.PackageCount > 20 {
		steps = append(steps,
			PlanStep{"buscar ciclos de dependencias"},
			PlanStep{"detectar paquetes demasiado grandes"},
		)
	}

	// Si hay archivos muy grandes, buscar hotspots
	if stats.LargeFiles > 10 {
		steps = append(steps,
			PlanStep{"detectar hotspots"},
		)
	}

	return Plan{Steps: steps}
}

func (p *AdvancedHeuristicPlanner) planExplainFile() Plan {
	return Plan{
		Steps: []PlanStep{
			{"extraer funciones"},
			{"extraer tipos"},
			{"extraer comentarios"},
			{"resumir archivo"},
		},
	}
}

func (p *AdvancedHeuristicPlanner) planCleanup() Plan {
	return Plan{
		Steps: []PlanStep{
			{"listar archivos"},
			{"limpiar imports"},
			{"formatear"},
			{"buscar código muerto"},
		},
	}
}

func (p *AdvancedHeuristicPlanner) planRefactor() Plan {
	return Plan{
		Steps: []PlanStep{
			{"buscar funciones largas"},
			{"buscar duplicación"},
			{"buscar nombres poco descriptivos"},
			{"buscar demasiados parámetros"},
		},
	}
}

func (p *AdvancedHeuristicPlanner) planSecurityAudit() Plan {
	return Plan{
		Steps: []PlanStep{
			{"buscar uso de exec.Command"},
			{"buscar rutas hardcodeadas"},
			{"buscar uso de md5 o sha1"},
			{"buscar errores ignorados"},
		},
	}
}

func (p *AdvancedHeuristicPlanner) planArchitecture() Plan {
	return Plan{
		Steps: []PlanStep{
			{"detectar dependencias"},
			{"buscar ciclos de dependencias"},
			{"buscar paquetes demasiado grandes"},
		},
	}
}

func (p *AdvancedHeuristicPlanner) planDocumentation() Plan {
	return Plan{
		Steps: []PlanStep{
			{"extraer comentarios"},
			{"extraer funciones"},
			{"generar documentación"},
		},
	}
}

//
// ============================
// ANÁLISIS DEL PROYECTO
// ============================
//

type ProjectStats struct {
	FileCount    int
	PackageCount int
	LargeFiles   int
}

func AnalyzeProjectSize() ProjectStats {
	// Aquí puedes usar tus tools reales
	// Por ahora devolvemos heurísticas simuladas
	return ProjectStats{
		FileCount:    350,
		PackageCount: 28,
		LargeFiles:   12,
	}
}

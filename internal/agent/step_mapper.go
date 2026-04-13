package agent

import "strings"

type StepMapper interface {
	Normalize(step string) string
}

type SemanticStepMapper struct{}

func NewSemanticStepMapper() *SemanticStepMapper {
	return &SemanticStepMapper{}
}

func (m *SemanticStepMapper) Normalize(step string) string {
	s := strings.ToLower(strings.TrimSpace(step))

	// === Listar archivos ===
	if containsAny(s, "listar archivos", "file tree", "ver archivos", "explorar archivos") {
		return "listar archivos"
	}

	// === Métricas ===
	if containsAny(s, "calcular métricas", "metrics", "estadísticas", "stats") {
		return "calcular métricas"
	}

	// === Dependencias ===
	if containsAny(s, "dependencias", "dependencies", "grafo de dependencias") {
		return "detectar dependencias"
	}

	// === Funciones largas ===
	if containsAny(s, "funciones largas", "long functions", "detectar funciones largas") {
		return "funciones largas"
	}

	// === Duplicación ===
	if containsAny(s, "duplicación", "duplicate code", "código duplicado") {
		return "duplicación"
	}

	// === Limpieza de imports ===
	if containsAny(s, "limpiar imports", "clean imports", "organizar imports") {
		return "limpiar imports"
	}

	// === Formatear código ===
	if containsAny(s, "formatear", "format code", "formato") {
		return "formatear"
	}

	// === Explicar archivo ===
	if containsAny(s, "explicar archivo", "explain file", "entender archivo") {
		return "explicar archivo"
	}

	// === Resumir archivo ===
	if containsAny(s, "resumir archivo", "summary", "resumen") {
		return "resumir archivo"
	}

	// === Extraer funciones ===
	if containsAny(s, "extraer funciones", "list functions", "obtener funciones") {
		return "extraer funciones"
	}

	// === Extraer tipos ===
	if containsAny(s, "extraer tipos", "list types", "obtener tipos") {
		return "extraer tipos"
	}

	// === Extraer comentarios ===
	if containsAny(s, "extraer comentarios", "list comments", "obtener comentarios") {
		return "extraer comentarios"
	}

	// === Código muerto ===
	if containsAny(s, "código muerto", "dead code") {
		return "dead code"
	}

	// Si no matchea nada, devolvemos el original
	return step
}

package agent

// Estas policies ahora mismo son simples wrappers de tu policy base.
// Luego puedes afinarlas con lógica específica por agente.

func NewAnalysisPolicy() AgentPolicy {
	// TODO: si tienes una policy específica de análisis, ponla aquí.
	// Por ahora reutilizamos la policy general.
	return NewDefaultPolicy()
}

func NewRefactorPolicy() AgentPolicy {
	// TODO: lógica específica de refactor (por ejemplo, priorizar tools de formato, extracción, etc.)
	return NewDefaultPolicy()
}

func NewDocsPolicy() AgentPolicy {
	// TODO: lógica específica de documentación (resumen, extracción de comentarios, etc.)
	return NewDefaultPolicy()
}

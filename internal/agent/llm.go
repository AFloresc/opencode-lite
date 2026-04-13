package agent

type LLMClient interface {
	// Dado un goal en lenguaje natural, devuelve una lista de pasos de alto nivel.
	ProposePlan(goal string) ([]string, error)
}

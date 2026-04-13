package agent

type DefaultPolicy struct{}

// Constructor
func NewDefaultPolicy() AgentPolicy {
	return &DefaultPolicy{}
}

// Decide es el fallback cuando no hay grounding directo.
// Aquí simplemente no hace nada y marca done=true.
func (p *DefaultPolicy) Decide(ctx *AgentContext) (string, map[string]interface{}, bool) {
	// No hay tool que ejecutar → terminamos
	return "", nil, true
}

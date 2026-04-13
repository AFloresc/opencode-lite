package agent

import "opencode-lite/internal/tools"

type MasterAgent struct {
	Agents     []SpecializedAgent
	Classifier *GoalClassifier
	Supervisor *Supervisor
}

// ------------------------------------------------------------
// Constructor final del MasterAgent
// ------------------------------------------------------------
func NewMasterAgent(projectID string, llm LLMClient) *MasterAgent {
	// Memoria cognitiva global del proyecto
	cogMem := NewCognitiveMemory(projectID)
	_ = cogMem.Load()

	// Crear agentes especializados con memoria compartida
	agents := []SpecializedAgent{
		NewAnalysisAgent(projectID, llm, cogMem),
		NewRefactorAgent(projectID, llm, cogMem),
		NewDocsAgent(projectID, llm, cogMem),
	}

	return &MasterAgent{
		Agents:     agents,
		Classifier: NewGoalClassifier(llm),
		Supervisor: NewSupervisor(llm, cogMem),
	}
}

// ------------------------------------------------------------
// Selección de agente por clasificación LLM
// ------------------------------------------------------------
func (m *MasterAgent) SelectAgent(goal string) SpecializedAgent {
	name := m.Classifier.Classify(goal)

	for _, a := range m.Agents {
		if a.Name() == name {
			return a
		}
	}

	// fallback: análisis
	return m.Agents[0]
}

// ------------------------------------------------------------
// Selección directa por nombre
// ------------------------------------------------------------
func (m *MasterAgent) getAgentByName(name string) SpecializedAgent {
	for _, a := range m.Agents {
		if a.Name() == name {
			return a
		}
	}
	return m.Agents[0]
}

// ------------------------------------------------------------
// Ejecución principal del MasterAgent
// ------------------------------------------------------------
func (m *MasterAgent) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	// 1. Seleccionar agente inicial
	agent := m.SelectAgent(goal)

	// 2. Obtener runtime del agente
	base, ok := agent.(*BaseSpecializedAgent)
	if !ok {
		return ctx
	}
	rt := base.runtime

	// 3. Supervisor analiza el goal con runtime + contexto
	decision := m.Supervisor.Analyze(goal, rt, &ctx)

	switch decision.Action {

	// --------------------------------------------------------
	// Pedir aclaración
	// --------------------------------------------------------
	case "clarify":
		ctx.LastResult = tools.ToolResult{Result: decision.Message}
		return ctx

	// --------------------------------------------------------
	// Dividir goal en sub‑goals
	// --------------------------------------------------------
	case "split":
		for _, sub := range decision.SubGoals {
			subCtx := m.Run(sub)
			ctx.History = append(ctx.History, subCtx.History...)
		}
		return ctx

	// --------------------------------------------------------
	// Replanificar con el mismo agente
	// --------------------------------------------------------
	case "replan":
		return agent.Run(goal, &ctx)

	// --------------------------------------------------------
	// Finalizar
	// --------------------------------------------------------
	case "finish":
		ctx.LastResult = tools.ToolResult{Result: decision.Message}
		return ctx

	// --------------------------------------------------------
	// Delegar a otro agente
	// --------------------------------------------------------
	case "delegate":
		agent = m.getAgentByName(decision.AgentName)
		return agent.Run(goal, &ctx)
	}

	return ctx
}

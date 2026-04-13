package agent

import "opencode-lite/internal/tools"

type MasterAgent struct {
	Agents     []SpecializedAgent
	Classifier *GoalClassifier
	Supervisor *Supervisor
}

func NewMasterAgent(projectID string, llm LLMClient) *MasterAgent {
	return &MasterAgent{
		Agents: []SpecializedAgent{
			NewAnalysisAgent(projectID, llm),
			NewRefactorAgent(projectID, llm),
			NewDocsAgent(projectID, llm),
			// aquí puedes añadir security, architecture, etc.
		},
		Classifier: NewGoalClassifier(llm),
		Supervisor: NewSupervisor(llm),
	}
}

func (m *MasterAgent) SelectAgent(goal string) SpecializedAgent {
	for _, a := range m.Agents {
		if a.CanHandle(goal) {
			return a
		}
	}
	// fallback: primer agente (o uno genérico)
	if len(m.Agents) > 0 {
		return m.Agents[0]
	}
	return nil
}

func (m *MasterAgent) Run(goal string) AgentContext {
	ctx := AgentContext{
		Goal:   goal,
		Memory: map[string]interface{}{},
	}

	// 1. Seleccionamos el agente inicial según el goal
	agent := m.SelectAgent(goal)

	// 2. Obtenemos su runtime (BaseSpecializedAgent → runtime)
	base, ok := agent.(*BaseSpecializedAgent)
	if !ok {
		// fallback improbable pero seguro
		return ctx
	}
	rt := base.runtime

	// 3. El Supervisor analiza el goal con acceso al runtime y al contexto
	decision := m.Supervisor.Analyze(goal, rt, &ctx)

	switch decision.Action {

	case "clarify":
		ctx.LastResult = tools.ToolResult{Result: decision.Message}
		return ctx

	case "split":
		for _, sub := range decision.SubGoals {
			subCtx := m.Run(sub)
			ctx.History = append(ctx.History, subCtx.History...)
		}
		return ctx

	case "replan":
		// Reutilizamos el mismo agente, pero con replanificación interna
		return agent.Run(goal, &ctx)

	case "finish":
		ctx.LastResult = tools.ToolResult{Result: decision.Message}
		return ctx

	case "delegate":
		agent = m.getAgentByName(decision.AgentName)
		return agent.Run(goal, &ctx)
	}

	return ctx
}

func (m *MasterAgent) getAgentByName(name string) SpecializedAgent {
	for _, a := range m.Agents {
		if a.Name() == name {
			return a
		}
	}
	return m.Agents[0]
}

package agent

type MasterAgent struct {
	Agents     []SpecializedAgent
	Classifier *GoalClassifier
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

	agent := m.SelectAgent(goal)
	if agent == nil {
		return ctx
	}

	return agent.Run(goal, &ctx)
}

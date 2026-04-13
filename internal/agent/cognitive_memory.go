package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CognitiveMemory struct {
	ProjectID string
	Data      map[string]interface{}
}

func NewCognitiveMemory(projectID string) *CognitiveMemory {
	return &CognitiveMemory{
		ProjectID: projectID,
		Data:      map[string]interface{}{},
	}
}

func (m *CognitiveMemory) path() string {
	return filepath.Join(".opencode", m.ProjectID, "cognitive_memory.json")
}

func (m *CognitiveMemory) Load() error {
	path := m.path()
	b, err := os.ReadFile(path)
	if err != nil {
		return nil // no existe → memoria vacía
	}
	return json.Unmarshal(b, &m.Data)
}

func (m *CognitiveMemory) Save() error {
	path := m.path()
	os.MkdirAll(filepath.Dir(path), 0755)
	b, _ := json.MarshalIndent(m.Data, "", "  ")
	return os.WriteFile(path, b, 0644)
}

func (m *CognitiveMemory) Remember(key string, value interface{}) {
	m.Data[key] = value
}

func (m *CognitiveMemory) Recall(key string) interface{} {
	return m.Data[key]
}

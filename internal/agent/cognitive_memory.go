package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
)

//
// ============================================================
//  CognitiveMemory
//  - memoria cognitiva persistente por proyecto
//  - almacena señales, patrones, últimos tools, últimos resultados
//  - usada por Supervisor, Grounder, Runtime, AOC, Metacognición
// ============================================================
//

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

//
// ============================================================
//  Load / Save
// ============================================================
//

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
	_ = os.MkdirAll(filepath.Dir(path), 0755)

	b, _ := json.MarshalIndent(m.Data, "", "  ")
	return os.WriteFile(path, b, 0644)
}

//
// ============================================================
//  API básica
// ============================================================
//

func (m *CognitiveMemory) Remember(key string, value interface{}) {
	m.Data[key] = value
}

func (m *CognitiveMemory) Recall(key string) interface{} {
	return m.Data[key]
}

func (m *CognitiveMemory) Reset(key string) {
	delete(m.Data, key)
}

func (m *CognitiveMemory) Increment(key string) int {
	v, ok := m.Data[key].(int)
	if !ok {
		m.Data[key] = 1
		return 1
	}
	m.Data[key] = v + 1
	return v + 1
}

//
// ============================================================
//  Helpers cognitivos de alto nivel
// ============================================================
//

// Última herramienta usada
func (m *CognitiveMemory) RememberLastTool(name string) {
	m.Remember("last_tool", name)
}

// Último resultado
func (m *CognitiveMemory) RememberLastResult(result interface{}) {
	m.Remember("last_result", result)
}

// Fallos repetidos
func (m *CognitiveMemory) RegisterFailure() int {
	return m.Increment("fail_count")
}

func (m *CognitiveMemory) ResetFailures() {
	m.Reset("fail_count")
}

// Éxitos repetidos
func (m *CognitiveMemory) RegisterSuccess() int {
	return m.Increment("success_count")
}

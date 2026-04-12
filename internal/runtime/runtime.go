package runtime

import (
	"encoding/json"
	"fmt"

	"opencode-lite/internal/providers"
	"opencode-lite/internal/tools"
)

type ChatRequest struct {
	Input string `json:"input"`
}

type ChatResponse struct {
	Message   string           `json:"message"`
	ToolCalls []tools.ToolCall `json:"tool_calls,omitempty"`
}

func HandleChatWithOllama(provider *providers.OllamaProvider, req ChatRequest) (ChatResponse, []tools.ToolResult, error) {

	messages := []providers.OllamaChatMessage{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: req.Input},
	}

	raw, err := provider.Chat(messages)
	if err != nil {
		return ChatResponse{}, nil, err
	}
	fmt.Println("RAW RESPONSE FROM MODEL:")
	fmt.Println(raw)

	parsed, ok := parseChatResponse(raw)
	if !ok {
		// No se pudo parsear como JSON válido (ni normal ni doble)
		return ChatResponse{
			Message: raw,
		}, nil, nil
	}

	// Ejecutar herramientas si las hay
	var results []tools.ToolResult
	for _, tc := range parsed.ToolCalls {
		res := tools.ExecuteTool(tc)
		results = append(results, res)
	}

	return parsed, results, nil
}

// parseChatResponse intenta:
// 1) parsear JSON normal
// 2) si falla, parsear JSON doble (cadena que contiene JSON)
func parseChatResponse(raw string) (ChatResponse, bool) {
	var parsed ChatResponse

	// 1. Intentar JSON normal
	if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
		return parsed, true
	}

	// 2. Intentar JSON doble: raw es una cadena que contiene JSON
	var inner string
	if err := json.Unmarshal([]byte(raw), &inner); err == nil {
		if err := json.Unmarshal([]byte(inner), &parsed); err == nil {
			return parsed, true
		}
	}

	// Nada funcionó
	fmt.Println("no se pudo parsear JSON de la respuesta del modelo")
	return ChatResponse{}, false
}

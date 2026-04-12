package runtime

import (
	"encoding/json"

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

	var parsed ChatResponse
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		// Si el modelo no devolvió JSON válido
		return ChatResponse{
			Message: raw,
		}, nil, nil
	}

	var results []tools.ToolResult
	for _, tc := range parsed.ToolCalls {
		res := tools.ExecuteTool(tc)
		results = append(results, res)
	}

	return parsed, results, nil
}

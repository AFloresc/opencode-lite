package runtime

import (
	"encoding/json"
	"fmt"
	"strings"

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

func HandleChatWithOllama(provider *providers.OllamaProvider, req ChatRequest, systemPrompt string) (ChatResponse, []tools.ToolResult, error) {

	// Historial inicial
	messages := []providers.OllamaChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: req.Input},
	}

	// 1. Primer turno: el modelo decide si usar herramientas
	raw, err := provider.Chat(messages)
	if err != nil {
		return ChatResponse{}, nil, err
	}

	fmt.Println("RAW RESPONSE FROM MODEL:")
	fmt.Println(raw)

	parsed, ok := parseChatResponse(raw)
	if !ok {
		// No es JSON válido → devolvemos texto
		return ChatResponse{Message: raw}, nil, nil
	}

	// Si no hay tool_calls → respuesta final
	if len(parsed.ToolCalls) == 0 {
		return parsed, nil, nil
	}

	// 2. Ejecutar herramientas
	var toolResults []tools.ToolResult

	for _, tc := range parsed.ToolCalls {
		res := tools.ExecuteTool(tc)
		toolResults = append(toolResults, res)

		// Extraer path
		path := ""
		if p, ok := tc.Arguments["path"].(string); ok {
			path = p
		}

		// Inyectar contenido real del archivo al modelo
		messages = []providers.OllamaChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf(
				"Contenido real del archivo %s:\n\n%s\n\nAhora responde basándote en este contenido.",
				path,
				res.Result,
			)},
		}

	}

	// 3. Segundo turno: el modelo responde basándose en el contenido real
	raw2, err := provider.Chat(messages)
	if err != nil {
		return ChatResponse{}, toolResults, err
	}

	fmt.Println("RAW SECOND RESPONSE FROM MODEL:")
	fmt.Println(raw2)

	parsed2, ok := parseChatResponse(raw2)
	if !ok {
		return ChatResponse{Message: raw2}, toolResults, nil
	}

	// Si el segundo turno también trae tool_calls (como apply_patch), ejecútalas
	if len(parsed2.ToolCalls) > 0 {
		for _, tc := range parsed2.ToolCalls {
			res := tools.ExecuteTool(tc)
			toolResults = append(toolResults, res)
		}

		// Opcional: respuesta final simple tras ejecutar tools del segundo turno
		return ChatResponse{
			Message: "herramientas del segundo turno ejecutadas",
		}, toolResults, nil
	}

	return parsed2, toolResults, nil
}

//
// PARSER ROBUSTO: extrae el primer objeto JSON válido usando conteo de llaves
//

func parseChatResponse(raw string) (ChatResponse, bool) {
	var parsed ChatResponse

	start := -1
	braceCount := 0

	for i, r := range raw {
		if r == '{' {
			if start == -1 {
				start = i
			}
			braceCount++
		} else if r == '}' {
			if braceCount > 0 {
				braceCount--
				if braceCount == 0 && start != -1 {
					jsonText := raw[start : i+1]
					jsonText = strings.TrimSpace(jsonText)

					if err := json.Unmarshal([]byte(jsonText), &parsed); err == nil {
						return parsed, true
					}

					fmt.Println("no se pudo parsear JSON de la respuesta del modelo")
					fmt.Println(jsonText)
					return ChatResponse{}, false
				}
			}
		}
	}

	fmt.Println("no se pudo parsear JSON de la respuesta del modelo")
	return ChatResponse{}, false
}

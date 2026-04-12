package providers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type OllamaProvider struct {
	BaseURL string
	Model   string
}

type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type ollamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: baseURL,
		Model:   model,
	}
}

func (p *OllamaProvider) Chat(messages []OllamaChatMessage) (string, error) {

	reqBody := ollamaChatRequest{
		Model:    p.Model,
		Messages: messages,
		Stream:   false,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(p.BaseURL+"/api/chat", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var out ollamaChatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		// Si Ollama devuelve texto plano, lo devolvemos tal cual
		return string(body), nil
	}

	return out.Message.Content, nil
}

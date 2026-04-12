package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
}

type OllamaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: baseURL,
		Model:   model,
	}
}

func (p *OllamaProvider) Chat(messages []OllamaChatMessage) (string, error) {
	reqBody := OllamaChatRequest{
		Model:    p.Model,
		Messages: messages,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(p.BaseURL+"/chat/completions", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return out.Choices[0].Message.Content, nil
}

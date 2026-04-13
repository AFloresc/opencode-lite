package agent

import (
	"encoding/json"
	"errors"
	"strings"
)

type LLMClient interface {
	ProposePlan(goal string) ([]string, error)
	Complete(prompt string) (string, error)
}

type PromptLLMClient struct {
	Call func(prompt string) (string, error)
}

func NewPromptLLMClient(callFn func(prompt string) (string, error)) *PromptLLMClient {
	return &PromptLLMClient{Call: callFn}
}

func (c *PromptLLMClient) ProposePlan(goal string) ([]string, error) {
	if c.Call == nil {
		return nil, errors.New("LLM call function not set")
	}

	prompt := strings.Replace(LLMPlannerPrompt, "GOAL_HERE", goal, 1)

	raw, err := c.Call(prompt)
	if err != nil {
		return nil, err
	}

	var steps []string
	if err := json.Unmarshal([]byte(raw), &steps); err != nil {
		return nil, errors.New("LLM did not return valid JSON list")
	}

	return steps, nil
}

func (c *PromptLLMClient) Complete(prompt string) (string, error) {
	if c.Call == nil {
		return "", errors.New("LLM call function not set")
	}

	return c.Call(prompt)
}

package config

import (
	"encoding/json"
	"os"
)

type ModelConfig struct {
	Model string `json:"model"`
}

type ProviderConfig struct {
	BaseURL string                 `json:"baseURL"`
	Models  map[string]ModelConfig `json:"models"`
}

type Config struct {
	Providers    map[string]ProviderConfig `json:"providers"`
	DefaultModel string                    `json:"defaultModel"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

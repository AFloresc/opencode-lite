package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"opencode-lite/internal/config"
	"opencode-lite/internal/providers"
	"opencode-lite/internal/runtime"
)

func main() {
	cfg, err := config.Load("opencode.json")
	if err != nil {
		fmt.Println("error loading config:", err)
		return
	}

	// Cargamos el provider localtools → Ollama
	p := providers.NewOllamaProvider(
		cfg.Providers["localtools"].BaseURL,
		cfg.Providers["localtools"].Models["qwen-local"].Model,
	)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}

		resp, toolResults, err := runtime.HandleChatWithOllama(p, runtime.ChatRequest{Input: line})
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		fmt.Println(resp.Message)
		for _, tr := range toolResults {
			if tr.Error != "" {
				fmt.Println("Tool error:", tr.Error)
			} else {
				fmt.Println("Tool OK:", tr.Result)
			}
		}
	}
}

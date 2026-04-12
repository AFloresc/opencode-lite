package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"opencode-lite/internal/providers"
	"opencode-lite/internal/runtime"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Configura tu Ollama local
	provider := providers.NewOllamaProvider(
		"http://localhost:11434", // URL de Ollama
		"qwen2.5-coder",          // Modelo que estés usando
	)

	fmt.Println("OpenCode Lite + Ollama")
	fmt.Println("Escribe 'exit' para salir.")
	fmt.Println("--------------------------------------------------")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Saliendo...")
			return
		}

		// Construimos la request
		req := runtime.ChatRequest{
			Input: input,
		}

		// Llamamos al runtime con el system prompt
		resp, toolResults, err := runtime.HandleChatWithOllama(provider, req, runtime.SystemPrompt)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Mostrar resultados de herramientas (si los hubo)
		for _, tr := range toolResults {
			if tr.Error != "" {
				fmt.Printf("[Tool %s ERROR]: %s\n", tr.ToolName, tr.Error)
			} else {
				fmt.Printf("[Tool %s OK]\n", tr.ToolName)
			}
		}

		// Mostrar respuesta final del modelo
		if resp.Message != "" {
			fmt.Println(resp.Message)
		}
	}
}

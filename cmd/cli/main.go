package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"opencode-lite/internal/agent"
	"opencode-lite/internal/providers"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Configura tu proveedor Ollama
	provider := providers.NewOllamaProvider(
		"http://localhost:11434",
		"qwen2.5-coder",
	)

	// Adaptamos el provider a tu LLMClient
	llm := agent.NewPromptLLMClient(func(prompt string) (string, error) {
		return provider.Chat([]providers.OllamaChatMessage{
			{Role: "user", Content: prompt},
		})
	})

	// Creamos el MasterAgent (multi‑agente)
	master := agent.NewMasterAgent("workspace", llm)

	fmt.Println("OpenCode Lite — Multi‑Agente")
	fmt.Println("Escribe 'exit' para salir.")
	fmt.Println("--------------------------------------------------")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Saliendo…")
			return
		}

		// Ejecutamos el goal con el MasterAgent
		ctx := master.Run(input)

		// Mostramos el historial de pasos ejecutados
		fmt.Println("--------------------------------------------------")
		fmt.Println("Historial de ejecución:")
		for _, step := range ctx.History {
			fmt.Printf("- [%s] %v\n", step.Action, step.Input)
		}

		// Mostramos el último resultado
		fmt.Println("--------------------------------------------------")
		fmt.Println("Resultado final:")
		fmt.Printf("%+v\n", ctx.LastResult)
		fmt.Println("--------------------------------------------------")
	}
}

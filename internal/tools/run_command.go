package tools

// run_command: ejecuta comandos sandboxed dentro del runtime del agente
func runCommandTool(args map[string]interface{}) ToolResult {
	return runCommandDispatch(args)
}

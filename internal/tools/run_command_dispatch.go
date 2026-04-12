package tools

import "strings"

var runCommandHandlers = map[string]func(string, []string) ToolResult{
	// básicos
	"count_lines":   runCmdCountLines,
	"file_size":     runCmdFileSize,
	"validate_json": runCmdValidateJSON,
	"echo":          runCmdEcho,
	"word_count":    runCmdWordCount,
	"char_count":    runCmdCharCount,
	"sha256":        runCmdSHA256,
	"list_dir":      runCmdListDir,
	"head":          runCmdHead,
	"tail":          runCmdTail,
	"search":        runCmdSearch,
	"now":           runCmdNow,

	// análisis de código
	"count_funcs":     runCmdCountFuncs,
	"count_imports":   runCmdCountImports,
	"find_structs":    runCmdFindStructs,
	"find_interfaces": runCmdFindInterfaces,

	// análisis de proyecto
	"project_stats": runCmdProjectStats,
	"largest_files": runCmdLargestFiles,
	"file_tree":     runCmdFileTree,

	// inteligencia
	"detect_language":  runCmdDetectLanguage,
	"summarize_file":   runCmdSummarizeFile,
	"extract_comments": runCmdExtractComments,
}

func runCommandDispatch(args map[string]interface{}) ToolResult {
	cmdRaw, ok := args["cmd"]
	if !ok {
		return ToolResult{"run_command", nil, "falta argumento obligatorio: cmd"}
	}

	cmd, ok := cmdRaw.(string)
	if !ok {
		return ToolResult{"run_command", nil, "el argumento 'cmd' debe ser string"}
	}

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ToolResult{"run_command", nil, "comando vacío"}
	}

	main := parts[0]
	argsList := parts[1:]

	// Tabla de comandos → función
	if handler, ok := runCommandHandlers[main]; ok {
		return handler(main, argsList)
	}

	return ToolResult{"run_command", nil, "comando no permitido o desconocido: " + main}
}

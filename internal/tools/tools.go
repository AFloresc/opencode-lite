package tools

import (
	"fmt"
)

//
// Estructuras base
//

// ToolCall representa una llamada a herramienta generada por el modelo
type ToolCallOLD struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult representa el resultado de ejecutar una herramienta
type ToolResultOLD struct {
	ToolName string      `json:"tool_name"`
	Result   interface{} `json:"result"`
	Error    string      `json:"error,omitempty"`
}

//
// Registro de herramientas disponibles
//

var toolRegistryOLD = map[string]func(map[string]interface{}) ToolResult{
	"read_file":              readFileTool,
	"write_file":             writeFileTool,
	"apply_patch":            applyPatchTool,
	"apply_patch_fuzzy":      applyPatchFuzzyTool,
	"apply_patch_auto":       applyPatchAutoTool,
	"list_files":             listFilesTool,
	"search_in_file":         searchInFileTool,
	"grep":                   grepTool,
	"delete_file":            deleteFileTool,
	"rename_file":            renameFileTool,
	"copy_file":              copyFileTool,
	"move_file":              moveFileTool,
	"create_file":            createFileTool,
	"file_exists":            fileExistsTool,
	"read_dir":               readDirTool,
	"append_file":            appendFileTool,
	"truncate_file":          truncateFileTool,
	"stat_file":              statFileTool,
	"touch_file":             touchFileTool,
	"search_replace":         searchReplaceTool,
	"diff_files":             diffFilesTool,
	"zip_dir":                zipDirTool,
	"unzip":                  unzipTool,
	"search_regex":           searchRegexTool,
	"apply_patch_structured": applyPatchStructuredTool,
	"format_code":            formatCodeTool,
	"search_regex_multi":     searchRegexMultiTool,
	"lint_code":              lintCodeTool,
	"run_command":            runCommandTool,
}

//
// Función principal para ejecutar herramientas
//

func ExecuteTool(call ToolCall) ToolResult {
	toolFunc, exists := ToolRegistry[call.Name]
	if !exists {
		return ToolResult{
			ToolName: call.Name,
			Error:    fmt.Sprintf("herramienta desconocida: %s", call.Name),
		}
	}

	return toolFunc(call.Arguments)
}

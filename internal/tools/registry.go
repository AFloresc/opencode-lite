package tools

var toolRegistry = map[string]func(map[string]interface{}) ToolResult{
	// Filesystem
	"read_file":     readFileTool,
	"write_file":    writeFileTool,
	"create_file":   createFileTool,
	"list_files":    listFilesTool,
	"delete_file":   deleteFileTool,
	"rename_file":   renameFileTool,
	"copy_file":     copyFileTool,
	"move_file":     moveFileTool,
	"file_exists":   fileExistsTool,
	"read_dir":      readDirTool,
	"stat_file":     statFileTool,
	"touch_file":    touchFileTool,
	"diff_files":    diffFilesTool,
	"append_file":   appendFileTool,
	"truncate_file": truncateFileTool,

	// Seacrh
	"search_in_file":     searchInFileTool,
	"grep":               grepTool,
	"search_replace":     searchReplaceTool,
	"search_regex":       searchRegexTool,
	"search_regex_multi": searchRegexMultiTool,

	// Patches
	"apply_patch":            applyPatchTool,
	"apply_patch_fuzzy":      applyPatchFuzzyTool,
	"apply_patch_auto":       applyPatchAutoTool,
	"apply_patch_structured": applyPatchStructuredTool,

	// Formating
	"format_code": formatCodeTool,

	// Compression
	"zip_dir": zipDirTool,
	"unzip":   unzipTool,

	// Lint
	"lint_code": lintCodeTool,

	// run_command
	"run_command": runCommandTool,

	// Project analysis

	// Semantic
}

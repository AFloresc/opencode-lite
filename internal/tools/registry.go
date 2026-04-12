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

	// Code analysis
	"count_funcs":     runCommandTool,
	"count_imports":   runCommandTool,
	"find_structs":    runCommandTool,
	"find_interfaces": runCommandTool,

	//Project analysis
	"project_stats": runCommandTool,
	"largest_files": runCommandTool,
	"file_tree":     runCommandTool,

	// Inteligence / semantic
	"detect_language":  runCommandTool,
	"summarize_file":   runCommandTool,
	"extract_comments": runCommandTool,

	// Semantic tools
	"extract_functions":      extractFunctionsTool,
	"extract_types":          extractTypesTool,
	"extract_comments_block": extractCommentsBlockTool,
	"semantic_index":         semanticIndexTool,

	// Refactor tools
	"refactor_rename_symbol": refactorRenameSymbolTool,
	"refactor_move_file":     refactorMoveFileTool,
	"refactor_split_file":    refactorSplitFileTool,
	"refactor_merge_files":   refactorMergeFilesTool,

	// Analysis tools
	"analysis_dependencies": analysisDependenciesTool,
	"analysis_cyclomatic":   analysisCyclomaticTool,
	"analysis_dead_code":    analysisDeadCodeTool,
	"analysis_metrics":      analysisMetricsTool,
}

package tools

import (
	"path/filepath"
	"regexp"
	"strings"
)

func runCmdDetectLanguage(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: detect_language <archivo>"}
	}
	ext := strings.ToLower(filepath.Ext(argsList[0]))
	lang := "unknown"
	switch ext {
	case ".go":
		lang = "go"
	case ".json":
		lang = "json"
	case ".yaml", ".yml":
		lang = "yaml"
	case ".md":
		lang = "markdown"
	case ".txt":
		lang = "text"
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"lang":    lang,
	}, ""}

}

func runCmdSummarizeFile(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: summarize_file <archivo>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	lines := strings.Split(content, "\n")
	summary := lines
	if len(summary) > 5 {
		summary = summary[:5]
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"summary": summary,
	}, ""}
}

func runCmdExtractComments(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: extract_comments <archivo.go>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	re := regexp.MustCompile(`(?m)//.*$`)
	comments := re.FindAllString(content, -1)
	return ToolResult{"run_command", map[string]interface{}{
		"command":  main,
		"comments": comments,
	}, ""}
}

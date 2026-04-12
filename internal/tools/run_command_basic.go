package tools

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func runCmdCountLines(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: count_lines <archivo>"}
	}

	path := filepath.Join("workspace", args[0])
	data, err := os.ReadFile(path)
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}

	lines := strings.Count(string(data), "\n") + 1

	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"lines":   lines,
	}, ""}
}

func runCmdFileSize(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: file_size <archivo>"}
	}
	path := filepath.Join("workspace", args[0])
	info, err := os.Stat(path)
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"size":    info.Size(),
	}, ""}
}

func runCmdValidateJSON(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: validate_json <archivo>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	var js interface{}
	err = json.Unmarshal([]byte(content), &js)
	if err != nil {
		return ToolResult{"run_command", map[string]interface{}{
			"command": main,
			"valid":   false,
			"error":   err.Error(),
		}, ""}
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"valid":   true,
	}, ""}
}

func runCmdEcho(main string, args []string) ToolResult {
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"output":  strings.Join(args, " "),
	}, ""}
}

func runCmdWordCount(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: word_count <archivo>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	words := len(strings.Fields(content))
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"words":   words,
	}, ""}

}

func runCmdCharCount(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: char_count <archivo>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"chars":   len(content),
	}, ""}
}

func runCmdSHA256(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: sha256 <archivo>"}
	}
	path := filepath.Join("workspace", args[0])
	data, err := os.ReadFile(path)
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	sum := sha256.Sum256(data)
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"sha256":  fmt.Sprintf("%x", sum),
	}, ""}
}

func runCmdListDir(main string, args []string) ToolResult {
	if len(args) != 1 {
		return ToolResult{"run_command", nil, "uso: list_dir <directorio>"}
	}
	path := filepath.Join("workspace", args[0])
	entries, err := os.ReadDir(path)
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"entries": names,
	}, ""}
}

func runCmdHead(main string, args []string) ToolResult {
	if len(args) != 2 {
		return ToolResult{"run_command", nil, "uso: head <archivo> <n>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	n, _ := strconv.Atoi(args[1])
	lines := strings.Split(content, "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"lines":   lines[:n],
	}, ""}
}

func runCmdTail(main string, args []string) ToolResult {
	if len(args) != 2 {
		return ToolResult{"run_command", nil, "uso: tail <archivo> <n>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	n, _ := strconv.Atoi(args[1])
	lines := strings.Split(content, "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"lines":   lines[len(lines)-n:],
	}, ""}
}

func runCmdSearch(main string, args []string) ToolResult {
	if len(args) < 2 {
		return ToolResult{"run_command", nil, "uso: search <archivo> <texto>"}
	}
	content, err := readFile(args[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	query := strings.Join(args[1:], " ")
	count := strings.Count(content, query)
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"query":   query,
		"count":   count,
	}, ""}
}

func runCmdNow(main string, args []string) ToolResult {
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"time":    time.Now().Format(time.RFC3339),
	}, ""}
}

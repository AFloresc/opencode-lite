package tools

import "regexp"

func runCmdCountFuncs(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: count_funcs <archivo.go>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	re := regexp.MustCompile(`(?m)^func\s+`)
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"count":   len(re.FindAllStringIndex(content, -1)),
	}, ""}
}

func runCmdCountImports(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: count_imports <archivo.go>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	re := regexp.MustCompile(`"([^"]+)"`)
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"imports": len(re.FindAllStringSubmatch(content, -1)),
	}, ""}
}

func runCmdFindStructs(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: find_structs <archivo.go>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	re := regexp.MustCompile(`type\s+([A-Za-z0-9_]+)\s+struct`)
	matches := re.FindAllStringSubmatch(content, -1)
	var names []string
	for _, m := range matches {
		names = append(names, m[1])
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command": main,
		"structs": names,
	}, ""}
}

func runCmdFindInterfaces(main string, argsList []string) ToolResult {
	if len(argsList) != 1 {
		return ToolResult{"run_command", nil, "uso: find_interfaces <archivo.go>"}
	}
	content, err := readFile(argsList[0])
	if err != nil {
		return ToolResult{"run_command", nil, err.Error()}
	}
	re := regexp.MustCompile(`type\s+([A-Za-z0-9_]+)\s+interface`)
	matches := re.FindAllStringSubmatch(content, -1)
	var names []string
	for _, m := range matches {
		names = append(names, m[1])
	}
	return ToolResult{"run_command", map[string]interface{}{
		"command":    main,
		"interfaces": names,
	}, ""}
}

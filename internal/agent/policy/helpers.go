package agent

import "strings"

func containsAny(s string, words ...string) bool {
	s = strings.ToLower(s)
	for _, w := range words {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}

func extractPattern(goal string) string {
	parts := strings.Split(goal, "\"")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ".*"
}

func extractFile(goal string) string {
	words := strings.Fields(goal)
	for _, w := range words {
		if strings.HasSuffix(w, ".go") {
			return w
		}
	}
	return ""
}

func extractFiles(goal string) []string {
	files := []string{}
	words := strings.Fields(goal)
	for _, w := range words {
		if strings.HasSuffix(w, ".go") {
			files = append(files, w)
		}
	}
	return files
}

func extractRename(goal string) (string, string) {
	words := strings.Fields(goal)
	var old, new string
	for i, w := range words {
		if w == "a" || w == "to" {
			if i > 0 {
				old = words[i-1]
			}
			if i < len(words)-1 {
				new = words[i+1]
			}
		}
	}
	return old, new
}

func extractNewPath(goal string) string {
	words := strings.Fields(goal)
	for i, w := range words {
		if w == "a" || w == "to" {
			if i < len(words)-1 {
				return words[i+1]
			}
		}
	}
	return ""
}

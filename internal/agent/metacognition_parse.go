package agent

import (
	"regexp"
	"strconv"
)

func extractConfidence(s string) float64 {
	re := regexp.MustCompile(`([0-1]\.\d+)`)
	m := re.FindString(s)
	if m == "" {
		return 0.5
	}
	f, _ := strconv.ParseFloat(m, 64)
	return f
}

func extractAdvice(s string) string {
	return s
}

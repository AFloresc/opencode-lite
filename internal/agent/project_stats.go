package agent

type ProjectStats struct {
	FileCount    int
	PackageCount int
	LargeFiles   int
}

func AnalyzeProjectSize() ProjectStats {
	// Aquí puedes usar tus tools reales más adelante.
	// Por ahora devolvemos valores simulados.
	return ProjectStats{
		FileCount:    350,
		PackageCount: 28,
		LargeFiles:   12,
	}
}

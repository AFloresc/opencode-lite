package agent

type Planner interface {
	Plan(goal string) []string
}

type SimplePlanner struct{}

func (p SimplePlanner) Plan(goal string) []string {
	if goal == "analizar proyecto" {
		return []string{
			"project_stats",
			"largest_files",
			"file_tree",
		}
	}
	return []string{goal}
}

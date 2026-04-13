package policy

import "opencode-lite/internal/agent"

type RuleBasedPolicy struct {
	rules []Rule
}

func NewRuleBasedPolicy() RuleBasedPolicy {

	return RuleBasedPolicy{
		rules: loadAllRules(),
	}
}

func (p RuleBasedPolicy) Decide(ctx *agent.AgentContext) (string, map[string]interface{}, bool) {
	goal := ctx.Goal

	for _, rule := range p.rules {
		if rule.Match(goal) {
			return rule.Apply(goal)
		}
	}

	return "", nil, true
}

type Rule struct {
	Match func(goal string) bool
	Apply func(goal string) (string, map[string]interface{}, bool)
}

func loadAllRules() []Rule {
	rules := []Rule{}
	rules = append(rules, searchRules()...)
	rules = append(rules, metricsRules()...)
	rules = append(rules, refactorRules()...)
	rules = append(rules, semanticRules()...)
	rules = append(rules, qualityRules()...)
	// Aquí puedes añadir más módulos: securityRules(), architectureRules(), etc.
	return rules
}

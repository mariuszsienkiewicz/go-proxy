package redirect

import (
	"go-proxy/modules/config"
)

type HashRule struct {
	Rule        config.Rule
	TargetGroup string
}

var (
	HashRules map[string]HashRule
)

func init() {
	HashRules = make(map[string]HashRule)
}

func BuildHashRules() {
	for _, rule := range config.Config.Proxy.Rules {
		if rule.Hash != "" {
			HashRules[rule.Hash] = HashRule{
				Rule:        rule,
				TargetGroup: rule.Target,
			}
		}
	}
}

func FindHashRule(hash string) (HashRule, bool) {
	rule, ok := HashRules[hash]
	return rule, ok
}

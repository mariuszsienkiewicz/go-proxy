package redirect

import (
	"proxy/modules/config"
)

type HashRule struct {
	Rule   config.Rule
	Target config.Server
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
				Rule:   rule,
				Target: ServerMap[rule.Target],
			}
		}
	}
}

func FindHashRule(hash string) (HashRule, bool) {
	rule, ok := HashRules[hash]
	return rule, ok
}

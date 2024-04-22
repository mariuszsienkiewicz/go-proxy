package redirect

import (
	"proxy/modules/config"
	"regexp"
)

var (
	RegexRules []RegexRule
)

type RegexRule struct {
	Rule        config.Rule
	Regex       string
	Regexp      *regexp.Regexp
	TargetGroup string
}

func (regexRule *RegexRule) Match(text string) bool {
	return regexRule.Regexp.MatchString(text)
}

func (regexRule *RegexRule) compile() {
	compiled, _ := regexp.Compile(regexRule.Regex)
	regexRule.Regexp = compiled
}

func BuildRegexRules() {
	for _, rule := range config.Config.Proxy.Rules {
		if rule.Regex != "" {
			r := RegexRule{
				Rule:        rule,
				Regex:       rule.Regex,
				TargetGroup: rule.Target,
			}
			r.compile()
			RegexRules = append(RegexRules, r)
		}
	}
}

func FindRegexRule(query string) (RegexRule, bool) {
	for _, regexRule := range RegexRules {
		if regexRule.Match(query) {
			return regexRule, true
		}
	}

	return RegexRule{}, false
}

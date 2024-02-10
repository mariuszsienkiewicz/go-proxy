package redirect

import (
	"proxy/modules/config"
	"regexp"
)

var (
	RegexRules []RegexRule
)

type RegexRule struct {
	Rule   config.Rule
	Regex  string
	Regexp *regexp.Regexp
	Target config.Server
}

func BuildRegexRules() {
	for _, rule := range config.Config.Proxy.Rules {
		r := RegexRule{
			Rule:  rule,
			Regex: rule.Regex,
		}
		r.compile()
		RegexRules = append(RegexRules, r)
	}
}

func (regexRule *RegexRule) compile() {
	compiled, _ := regexp.Compile(regexRule.Regex)
	regexRule.Regexp = compiled
}

func (regexRule *RegexRule) match(text string) bool {
	return regexRule.Regexp.MatchString(text)
}

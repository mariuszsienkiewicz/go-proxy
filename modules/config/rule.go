package config

import (
	"errors"
	"fmt"
)

type Rule struct {
	Name   string `yaml:"name"`
	Hash   string `yaml:"hash_rule,omitempty"`
	Regex  string `yaml:"regex_rule,omitempty"`
	Target string `yaml:"target_id"`
}

func ValidateRuleConfiguration() []error {
	errs := make([]error, 0)
	for i, rule := range Config.Proxy.Rules {
		if rule.Hash == "" && rule.Regex == "" {
			errs = append(errs, errors.New(fmt.Sprintf("[RULE %v ERROR] (%v): regex_rule or hash_rule must be specified", i+1, rule.Name)))
		}
	}

	if errs == nil || len(errs) == 0 {
		return nil
	}

	return errs
}

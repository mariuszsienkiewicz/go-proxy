package config

type Rule struct {
	Name   string `yaml:"name"`
	Regex  string `yaml:"regex_rule"`
	Target string `yaml:"target_id"`
}

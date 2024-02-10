package config

type DbUser struct {
	Target   string `yaml:"target"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

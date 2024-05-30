package config

import (
	"go-proxy/modules/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	// Config - in memory config representation of the yaml configuration file
	Config Configuration
)

// Configuration config wrapper, represents the whole yaml configuration file
type Configuration struct {
	Proxy ProxyConfig `yaml:"proxy"`
}

// ProxyConfig proxy related config
type ProxyConfig struct {
	Basics        Basics        `yaml:"basics"`
	Cache         Cache         `yaml:"cache,omitempty"`
	ServerGroups  []ServerGroup `yaml:"server_groups"`
	Servers       []Server      `yaml:"servers"`
	DbUsers       []DbUser      `yaml:"db_users"`
	Access        Access        `yaml:"access"`
	Rules         []Rule        `yaml:"rules"`
	DefaultServer *Server
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Proxy: ProxyConfig{
			Cache: GetDefaultCache(),
		},
	}
}

// LoadConfig loads the configuration to memory and verifies correctness of the configuration file
func LoadConfig(configPath string) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Logger.Fatal("Error while reading configuration file", zap.Error(err))
	}

	Config = *NewConfiguration()

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Logger.Fatal("Error while parsing configuration file", zap.Error(err))
	}

	if errs := validate(); errs != nil {
		log.Logger.Fatal("Configuration file validation failed", zap.Errors("errors", errs))
	}
}

func validate() []error {
	var errs []error
	if err := ValidateBasicConfiguration(); err != nil {
		return append(errs, err)
	}
	if err := ValidateServerConfiguration(); err != nil {
		return append(errs, err)
	}
	if err := ValidateRuleConfiguration(); err != nil {
		return append(errs, err...)
	}
	if err := ValidateCacheConfiguration(); err != nil {
		return append(errs, err)
	}

	return nil
}

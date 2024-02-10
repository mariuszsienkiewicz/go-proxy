package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"proxy/modules/log"
)

var (
	Config Configuration
)

type Configuration struct {
	Proxy struct {
		Basics        Basics   `yaml:"basics"`
		Servers       []Server `yaml:"servers"`
		DbUsers       []DbUser `yaml:"db_users"`
		Access        Access   `yaml:"access"`
		Rules         []Rule   `yaml:"rules"`
		DefaultServer *Server
	} `yaml:"proxy"`
}

func LoadConfig() {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Logger.Fatalf("Error while reading configuration file: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Logger.Fatalf("Error while parsing configuration file: %v\n", err)
	}

	if err := validate(); err != nil {
		log.Logger.Fatal(err)
	}
}

func validate() error {
	if err := ValidateBasicConfiguration(); err != nil {
		return err
	}

	if err := ValidateServerConfiguration(); err != nil {
		return err
	}

	return nil
}

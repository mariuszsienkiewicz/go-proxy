package config

import "fmt"

// TODO - add a way to change connection pool settings

type Basics struct {
	Port uint16 `yaml:"port"`
	Host string `yaml:"host"`
}

func (basics *Basics) GetHostname() string {
	return fmt.Sprintf("%v:%v", basics.Host, basics.Port)
}

func ValidateBasicConfiguration() error {
	return nil
}

package config

import (
	"errors"
	"fmt"
)

type Server struct {
	Name     string `yaml:"name"`
	Id       string `yaml:"id"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Required bool   `yaml:"required,omitempty"`
	TestDb   string `yaml:"test_db,omitempty"`
	Default  bool   `yaml:"default,omitempty"`
}

func (server *Server) GetDsn() string {
	return fmt.Sprintf("%s:%d", server.Host, server.Port)
}

// GetUser
// TODO: Search for the user and then cache it (or something) so it won't be necessary to iterate over users table again
func (server *Server) GetUser(users []DbUser) *DbUser {
	for _, user := range users {
		if user.Target == server.Id {
			return &user
		}
	}
	return nil
}

func ValidateServerConfiguration() error {
	// setup default server in configuration
	for _, server := range Config.Proxy.Servers {
		if server.Default {
			Config.Proxy.DefaultServer = &server
		}
	}

	// if default server was not in the config
	if Config.Proxy.DefaultServer == nil {
		// set automatically default server if there is only one server in configuration
		if len(Config.Proxy.Servers) == 1 {
			Config.Proxy.DefaultServer = &Config.Proxy.Servers[0]
		} else {
			return errors.New("none of the servers has 'default' property set as true")
		}
	}

	return nil
}

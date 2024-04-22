package config

import (
	"errors"
	"fmt"
)

type Server struct {
	Name        string `yaml:"name"`
	Id          string `yaml:"id"`
	Host        string `yaml:"host"`
	Port        uint16 `yaml:"port"`
	Required    bool   `yaml:"required,omitempty"`
	TestDb      string `yaml:"test_db,omitempty"`
	Default     bool   `yaml:"default,omitempty"`
	ServerGroup string `yaml:"server_group"`
}

var ErrNotFound = errors.New("user not found")

func (server *Server) GetDsn() string {
	return fmt.Sprintf("%s:%d", server.Host, server.Port)
}

func (server *Server) GetUser() (DbUser, error) {
	for _, user := range Config.Proxy.DbUsers {
		if user.Target == server.Id {
			return user, nil
		}
	}
	return DbUser{}, ErrNotFound
}

func ValidateServerConfiguration() error {
	// setup default db in configuration
	for _, server := range Config.Proxy.Servers {
		if server.Default {
			Config.Proxy.DefaultServer = &server
		}
	}

	// if default db was not in the config
	if Config.Proxy.DefaultServer == nil {
		// set automatically default db if there is only one db in configuration
		if len(Config.Proxy.Servers) == 1 {
			Config.Proxy.DefaultServer = &Config.Proxy.Servers[0]
		} else {
			return errors.New("none of the servers has 'default' property set as true")
		}
	}

	// TODO check if group exists

	return nil
}

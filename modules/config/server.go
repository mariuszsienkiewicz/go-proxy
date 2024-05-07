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
	// check if there is default db in configuration
	isAnyServerDefault := false
	for _, server := range Config.Proxy.Servers {
		if server.Default {
			isAnyServerDefault = true
			break
		}
	}

	if !isAnyServerDefault {
		return errors.New("no default server")
	}

	// TODO check if group exists

	return nil
}

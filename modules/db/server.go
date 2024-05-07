package db

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"proxy/modules/config"
)

type Server struct {
	Config      config.Server
	Credentials config.DbUser
	Status      Status
}

func LoadServers() error {
	for _, server := range config.Config.Proxy.Servers {
		g, groupExists := Groups[server.ServerGroup]
		if groupExists == false {
			return fmt.Errorf("server group %s does not exist", server.ServerGroup)
		}

		s, err := NewServer(server)
		if err != nil {
			return err
		}

		g.AddServer(s)
		DbPool.Servers[server.Id] = s
		if s.Config.Default {
			DbPool.DefaultServer = s
		}
	}

	return nil
}

func NewServer(server config.Server) (*Server, error) {
	user, err := server.GetUser()
	if err != nil {
		return &Server{}, err
	}

	return &Server{
		Config:      server,
		Credentials: user,
		Status:      SHUNNED, // by default, it has to be checked first
	}, nil
}

func (s Server) TestConnection() error {
	return TestConnection(s.Config.GetDsn(), s.Credentials.User, s.Credentials.Password, s.Config.TestDb)
}

func (s Server) Connect() (client.Conn, error) {
	return Connect(s.Config.GetDsn(), s.Credentials.User, s.Credentials.Password, s.Config.TestDb)
}

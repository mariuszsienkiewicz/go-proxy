package db

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"proxy/modules/config"
	"proxy/modules/log"
	"time"
)

type Server struct {
	Config      config.Server
	Credentials config.DbUser
	Status      Status
	Pool        *client.Pool
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

	pool := client.NewPool(log.Logger.Tracef, 100, 400, 5, fmt.Sprintf("%s:%d", server.Host, server.Port), user.User, user.Password, "")

	return &Server{
		Config:      server,
		Credentials: user,
		Status:      SHUNNED, // by default, it has to be checked first
		Pool:        pool,
	}, nil
}

func (s *Server) TestConnection() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer ctxCancel()

	conn, err := s.Pool.GetConn(ctx)
	if err != nil {
		return err
	}

	err = conn.Ping()

	// if conn can't connect to the database then drop this connection
	if err != nil {
		s.Pool.DropConn(conn)
		return err
	}

	// if connection was successful return to pool and return nil - everything is alright
	s.Pool.PutConn(conn)
	return nil
}

func (s *Server) Connect() (*client.Conn, error) {
	ctx := context.WithoutCancel(context.Background())

	conn, err := s.Pool.GetConn(ctx)
	if err != nil {
		return &client.Conn{}, err
	}

	return conn, nil
}

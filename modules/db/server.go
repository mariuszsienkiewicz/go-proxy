package db

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"go-proxy/modules/config"
	"go-proxy/modules/log"
	"time"
)

type Server struct {
	Config      config.Server
	Credentials config.DbUser
	Status      Status
	Pool        *client.Pool
}

func LoadServers(ctx context.Context) error {
	for _, server := range config.Config.Proxy.Servers {
		// Check if the context is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		g, groupExists := Groups[server.ServerGroup]
		if groupExists == false {
			return fmt.Errorf("server group %s does not exist", server.ServerGroup)
		}

		s, err := NewServer(ctx, server)
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

func NewServer(ctx context.Context, server config.Server) (*Server, error) {
	user, err := server.GetUser()
	if err != nil {
		return &Server{}, err
	}

	pool := client.NewPool(log.InfoWeak, 80, 150, 10, fmt.Sprintf("%s:%d", server.Host, server.Port), user.User, user.Password, "")

	// Create a server instance
	s := &Server{
		Config:      server,
		Credentials: user,
		Status:      SHUNNED, // by default, it has to be checked first
		Pool:        pool,
	}

	// Run a goroutine to close the pool when the context is done
	go func() {
		<-ctx.Done()
		s.Pool.Close()
	}()

	return s, nil
}

func (s *Server) TestConnection(ctx context.Context) error {
	ctxWithTimeout, ctxCancel := context.WithTimeout(ctx, 60000*time.Second)
	defer ctxCancel()

	conn, err := s.Pool.GetConn(ctxWithTimeout)
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

func (s *Server) Connect(ctx context.Context) (*client.Conn, error) {
	conn, err := s.Pool.GetConn(ctx)
	if err != nil {
		return &client.Conn{}, err
	}

	return conn, nil
}

package mysql

import (
	"github.com/go-mysql-org/go-mysql/client"
	"proxy/modules/config"
)

var (
	Connections map[string]client.Conn
)

type ConnectionContext struct {
	DbName     string
	Connection client.Conn
}

func init() {
	Connections = make(map[string]client.Conn)
}

func TestConnection(server config.Server, user config.DbUser) error {
	conn, err := client.Connect(server.GetDsn(), user.User, user.Password, server.TestDb)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func Connect(server config.Server, user config.DbUser) (client.Conn, error) {
	conn, ok := Connections[server.Id]
	if ok {
		return conn, nil
	}

	c, err := client.Connect(server.GetDsn(), user.User, user.Password, server.TestDb)
	if err != nil {
		return client.Conn{}, err
	}

	// add to connections
	Connections[server.Id] = *c

	return *c, nil
}

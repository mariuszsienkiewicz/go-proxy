package db

import (
	"github.com/go-mysql-org/go-mysql/client"
)

// TODO use Connection Pool

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

func TestConnection(addr string, user string, pass string, dbname string) error {
	conn, err := client.Connect(addr, user, pass, dbname)
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

func Connect(addr string, user string, pass string, dbname string) (client.Conn, error) {
	c, err := client.Connect(addr, user, pass, dbname)
	if err != nil {
		return client.Conn{}, err
	}

	return *c, nil
}

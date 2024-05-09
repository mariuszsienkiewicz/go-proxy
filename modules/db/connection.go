package db

import (
	"github.com/go-mysql-org/go-mysql/client"
)

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

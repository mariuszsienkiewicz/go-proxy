package cmd

import (
	"context"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/urfave/cli/v2"
	"net"
	"proxy/modules/config"
	"proxy/modules/db"
	"proxy/modules/log"
	"proxy/modules/mysql"
	"proxy/modules/redirect"
)

var Proxy = &cli.Command{
	Name:        "proxy",
	Usage:       "Start Proxy Service",
	Description: "",
	Action:      runProxy,
}

func runProxy(ctx *cli.Context) error {
	log.Logger.Tracef("Proxy command is running")

	log.Logger.Tracef("Setting up the command")
	setupError := setup()
	if setupError != nil {
		log.Logger.Fatal(setupError)
	}

	log.Logger.Tracef("Proxy is ready, serving")
	serve(ctx.Context)

	return nil
}

func setup() error {
	config.LoadConfig()
	redirect.BuildRules()

	log.Logger.Tracef("Initialization of db pools, groups and servers")
	err := db.Init()
	if err != nil {
		return err
	}

	return nil
}

func serve(ctx context.Context) {
	log.Logger.Infof("Proxy is running on: %v", config.Config.Proxy.Basics.GetHostname())

	// create TCP listener
	l, err := net.Listen("tcp", config.Config.Proxy.Basics.GetHostname())
	log.Logger.Infof("Listening on: %v", config.Config.Proxy.Basics.GetHostname())
	if err != nil {
		log.Logger.Fatal(err)
	}

	// close listener on function exit
	defer func() {
		if err := l.Close(); err != nil {
			log.Logger.Errorf("Error closing listener: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Logger.Tracef("Context canceled, shutting down the listener")
				return
			default:
				c, err := l.Accept()
				if err != nil {
					log.Logger.Errorf("Error accepting connection: %v", err)
					continue
				}

				go handleConnection(ctx, c)
			}
		}
	}()

	<-ctx.Done()
}

func handleConnection(ctx context.Context, c net.Conn) {
	log.Logger.Infof("Handle connection: %v", config.Config.Proxy.Basics.GetHostname())
	conn, err := server.NewConn(c, config.Config.Proxy.Access.User, config.Config.Proxy.Access.Password, &mysql.ProxyHandler{})
	if err != nil {
		log.Logger.Fatal(err)
	}

	// TODO there is a huge bug here, c.Close() will be executed when HandleCommand finishes
	for {
		select {
		case <-ctx.Done():
			log.Logger.Tracef("Closing connection with the client")
			err := c.Close()
			if err != nil {
				log.Logger.Error(err)
			}
			return
		default:
			if err := conn.HandleCommand(); err != nil {
				log.Logger.Error(err)
			}
		}
	}
}

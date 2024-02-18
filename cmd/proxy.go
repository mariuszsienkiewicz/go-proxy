package cmd

import (
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/urfave/cli/v2"
	"net"
	"proxy/modules/config"
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
	log.Logger.Tracef("Setup completed")

	log.Logger.Tracef("Proxy command is ready, serving")
	serve()

	return nil
}

func setup() error {
	config.LoadConfig()
	redirect.BuildRules()

	// check if it's possible to connect with the databases
	log.Logger.Tracef("Checking servers: %v", config.Config.Proxy.Servers)
	for _, dbServer := range config.Config.Proxy.Servers {
		if dbServer.Required {
			log.Logger.Tracef("Server %v is required, checking connectivity via %v", dbServer.Id, dbServer.GetDsn())
			if err := mysql.TestConnection(dbServer, *dbServer.GetUser(config.Config.Proxy.DbUsers)); err != nil {
				return err
			}
			log.Logger.Tracef("Server %v is recheable", dbServer.Id)
		}
	}

	return nil
}

func serve() {
	log.Logger.Infof("Proxy is running on: %v", config.Config.Proxy.Basics.GetHostname())

	l, err := net.Listen("tcp", config.Config.Proxy.Basics.GetHostname())
	log.Logger.Infof("Listening on: %v", config.Config.Proxy.Basics.GetHostname())
	if err != nil {
		log.Logger.Fatal(err)
	}

	for {
		c, err := l.Accept()
		log.Logger.Tracef("New connection accepted on: %v", config.Config.Proxy.Basics.GetHostname())
		if err != nil {
			log.Logger.Fatal(err)
		}

		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	log.Logger.Infof("Handle connection: %v", config.Config.Proxy.Basics.GetHostname())
	conn, err := server.NewConn(c, config.Config.Proxy.Access.User, config.Config.Proxy.Access.Password, &mysql.ProxyHandler{}) // TODO: accept only access user (from config.yml)
	if err != nil {
		log.Logger.Fatal(err)
	}

	for {
		if err := conn.HandleCommand(); err != nil {
			log.Logger.Error(err)
		}
	}
}

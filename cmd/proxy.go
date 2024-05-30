package cmd

import (
	"context"
	"errors"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go-proxy/modules/cache"
	"go-proxy/modules/config"
	"go-proxy/modules/db"
	"go-proxy/modules/log"
	"go-proxy/modules/proxy"
	"go-proxy/modules/redirect"
	"go.uber.org/zap"
	"net"
)

var Proxy = &cli.Command{
	Name:        "proxy",
	Usage:       "Start Proxy Service",
	Description: "",
	Action:      runProxy,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
	},
}

func runProxy(ctx *cli.Context) error {
	log.Logger.Info("Proxy command is running")

	log.Logger.Debug("Setting up the command")
	setupError := setup(ctx)
	if setupError != nil {
		log.Logger.Fatal("Setup error", zap.Error(setupError))
	}

	log.Logger.Info("Monitoring starting up...")
	db.MonitorServers(ctx.Context)

	log.Logger.Info("Proxy is ready, serving")
	serve(ctx.Context)

	return nil
}

func setup(ctx *cli.Context) error {
	// check if config file was set
	configPath := ctx.String("config")
	if configPath == "" {
		return errors.New("config file path is required")
	}

	// load config first
	config.LoadConfig(configPath)

	// initialize cache based on configuration
	if err := cache.InitCache(); err != nil {
		return err
	}

	// build regex rules
	redirect.BuildRules()

	log.Logger.Debug("Initialization of db pools, groups and servers")
	err := db.Init(ctx.Context)
	if err != nil {
		return err
	}

	return nil
}

func serve(ctx context.Context) {
	// create TCP listener
	l, err := net.Listen("tcp", config.Config.Proxy.Basics.GetHostname())
	if err != nil {
		log.Logger.Fatal("Listener error", zap.Error(err))
	}

	log.Logger.Info("Listening", zap.String("addr", config.Config.Proxy.Basics.GetHostname()))
	// close listener on function exit
	defer func() {
		if err := l.Close(); err != nil {
			log.Logger.Error("Error closing listener", zap.Error(err))
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Logger.Info("Context canceled, shutting down the listener")
				return
			default:
				c, err := l.Accept()
				if err != nil {
					log.Logger.Warn("Error accepting connection", zap.Error(err))
					continue
				}

				connectionId := uuid.New().String()
				log.Logger.Info("Accepting connection", zap.String("connection id", connectionId))
				go handleConnection(ctx, c, connectionId)
			}
		}
	}()

	<-ctx.Done()
}

func handleConnection(ctx context.Context, c net.Conn, connectionId string) {
	handler := proxy.NewProxyHandler(ctx, connectionId)
	defer handler.ConnectionManager.ReturnConnectionsToPool()

	conn, err := server.NewConn(c, config.Config.Proxy.Access.User, config.Config.Proxy.Access.Password, handler)
	if err != nil {
		log.Logger.Warn("Error creating new connection with proxy db proxy", zap.Error(err))
		err := c.Close()
		if err != nil {
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			log.Logger.Debug("Closing connection with the client")
			err := c.Close()
			if err != nil {
				log.Logger.Warn("Error while closing the connection", zap.Error(err))
			}
			return
		default:
			if err := conn.HandleCommand(); err != nil {
				log.Logger.Info("Handling command stopped", zap.String("handler", handler.Id), zap.NamedError("reason", err))
				return
			}
		}
	}
}

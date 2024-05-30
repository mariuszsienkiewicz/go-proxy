package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"go-proxy/modules/log"
	"go.uber.org/zap"
)

func NewProxyApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Proxy"
	app.HelpName = "proxy"
	app.Usage = "TODO"
	app.Description = "Proxy is a MySQL proxy service giving you possibility to redirect queries"
	app.Version = "1.0.0"

	subCmdWithConfig := []*cli.Command{
		Proxy,
	}

	app.Commands = append(app.Commands, subCmdWithConfig...)

	return app
}

func RunProxyApp(ctx context.Context, app *cli.App, args ...string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := app.RunContext(ctx, args)
			if err == nil {
				return nil
			} else {
				log.Logger.Error("Run Proxy App failed", zap.Error(err))
			}
		}
	}
}

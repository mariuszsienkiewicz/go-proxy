package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

func NewProxyApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Proxy"
	app.HelpName = "proxy"
	app.Usage = "TODO"
	app.Description = "Proxy is a MySQL proxy service giving you possibility to redirect queries"
	app.Version = "1.0.0"

	subCmdWithConfig := []*cli.Command{
		CmdProxy,
	}

	app.Commands = append(app.Commands, subCmdWithConfig...)

	return app
}

func RunProxyApp(app *cli.App, args ...string) error {
	err := app.Run(args)
	if err == nil {
		return nil
	}
	if strings.HasPrefix(err.Error(), "flag provided but not defined:") {
		cli.OsExiter(1)
		return err
	}
	_, _ = fmt.Fprintf(app.ErrWriter, "Command error: %v\n", err)
	cli.OsExiter(1)
	return err
}

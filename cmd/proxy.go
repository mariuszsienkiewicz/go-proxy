package cmd

import (
	"github.com/urfave/cli/v2"
)

var CmdProxy = &cli.Command{
	Name:        "proxy",
	Usage:       "Start Proxy Service",
	Description: "",
	Action:      runProxy,
}

func runProxy(ctx *cli.Context) error {
	return nil
}

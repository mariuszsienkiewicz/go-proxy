package main

import (
	"os"
	"proxy/cmd"
	"proxy/modules/log"
)

func main() {
	log.SetLogger()

	log.Logger.Trace("Creating the app")
	app := cmd.NewProxyApp()
	log.Logger.Trace("Running the app")
	_ = cmd.RunProxyApp(app, os.Args...)
}

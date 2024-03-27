package main

import (
	"context"
	"os"
	"proxy/cmd"
	"proxy/modules/log"
)

func main() {
	// setup context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	// setup logger
	log.SetLogger()

	// setup app
	log.Logger.Trace("Creating the app")
	app := cmd.NewProxyApp()
	log.Logger.Trace("Running the app")

	// run the app, exit on failure
	if err := cmd.RunProxyApp(app, os.Args...); err != nil {
		log.Logger.Errorf("%s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"proxy/cmd"
	"proxy/modules/log"
	"syscall"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	// setup logger
	log.SetLogger()

	// setup app
	log.Logger.Trace("Creating the app")
	app := cmd.NewProxyApp()

	// setup signal watcher
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			select {
			case sig := <-signalChan:
				switch sig {
				case syscall.SIGHUP:
					// reloading config requires the app to restart
					log.Logger.Infof("Signal: SIGHUP received, reloading config.")
					cancel()
				case syscall.SIGTERM, syscall.SIGINT:
					log.Logger.Infof("Signal %s received, exiting.", sig.String())
					cancel()
					os.Exit(1)
				}
			}
		}
	}()

	// cancel context when main function finishes
	defer cancel()

	for {
		run(app)
	}
}

func getNewContext() (context.Context, context.CancelFunc) {
	log.Logger.Tracef("Creating new cancelable context")
	newCtx, cancelFunc := context.WithCancel(context.Background())

	return newCtx, cancelFunc
}

// run the app, exit on failure
func run(app *cli.App) {
	log.Logger.Trace("Running the app")
	ctx, cancel = getNewContext()
	if err := cmd.RunProxyApp(ctx, app, os.Args...); err != nil {
		log.Logger.Errorf("%s\n", err)
		os.Exit(1)
	}
}

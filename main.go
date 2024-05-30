package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"go-proxy/cmd"
	"go-proxy/modules/log"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	// setup logger
	log.SetLogger()

	// Start HTTP server for pprof
	// http://localhost:6060/debug/pprof/profile?seconds=30
	//go func() {
	//	log.Logger.Info("Starting pprof server on :6060")
	//	err := http.ListenAndServe("localhost:6060", nil)
	//	if err != nil {
	//		return
	//	}
	//}()

	// setup app
	log.Logger.Info("Creating the app")
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
					log.Logger.Info("Signal: SIGHUP received, reloading config.")
					cancel()
				case syscall.SIGTERM, syscall.SIGINT:
					log.Logger.Info("Exit signal received, exiting.", zap.String("signal", sig.String()))
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
	log.Logger.Debug("Creating new cancelable context")
	newCtx, cancelFunc := context.WithCancel(context.Background())

	return newCtx, cancelFunc
}

// run the app, exit on failure
func run(app *cli.App) {
	log.Logger.Debug("Running the app")
	ctx, cancel = getNewContext()
	if err := cmd.RunProxyApp(ctx, app, os.Args...); err != nil {
		log.Logger.Debug("Error while running the proxy app", zap.Error(err))
		os.Exit(1)
	}
}

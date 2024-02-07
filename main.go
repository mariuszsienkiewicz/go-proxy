package main

import (
	"os"
	"proxy/cmd"
)

func main() {
	app := cmd.NewProxyApp()
	_ = cmd.RunProxyApp(app, os.Args...)
}

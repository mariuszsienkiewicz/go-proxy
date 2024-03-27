package log

import (
	"github.com/withmandala/go-log"
	"os"
)

var (
	Logger *log.Logger
)

func SetLogger() {
	Logger = log.New(os.Stdout).WithDebug() // TODO make debug optional
}

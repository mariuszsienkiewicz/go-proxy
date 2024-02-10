package log

import (
	"github.com/withmandala/go-log"
	"os"
)

var (
	Logger *log.Logger
)

func SetLogger() {
	Logger = log.New(os.Stderr).WithDebug() // TODO make debug optional
}

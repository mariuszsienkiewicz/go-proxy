package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	Sugared *zap.SugaredLogger
)

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stdout",
		},
	}

	return zap.Must(config.Build())
}

func InfoWeak(format string, args ...interface{}) {
	Sugared.Infof(format, args...)
}

func SetLogger() {
	Logger = createLogger()
	Sugared = createLogger().Sugar()
}

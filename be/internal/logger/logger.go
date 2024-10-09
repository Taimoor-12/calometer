package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func initLogger() {
	// Configure logger
	logConf := zap.NewProductionConfig()
	logConf.EncoderConfig.LevelKey = "l"
	logConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	// Build the logger
	var err error
	logger, err = logConf.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	if logger == nil {
		initLogger()
		logger.Sync()
	}

	return logger
}

package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {
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
	return logger
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

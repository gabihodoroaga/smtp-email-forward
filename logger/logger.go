package logger

import (
		"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Cannot create the logger")
	}
	defer logger.Sync()
	Log = logger.Sugar()
}

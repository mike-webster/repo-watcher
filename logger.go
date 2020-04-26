package main

import (
	"time"

	"github.com/gin-gonic/gin"
	env "github.com/mike-webster/repo-watcher/env"
	"github.com/sirupsen/logrus"
)

func defaultLogger(ctx *gin.Context) *logrus.Logger {
	if ctx == nil {
		return newLogger()
	}

	var logger *logrus.Logger
	l, exists := ctx.Get("logger")
	if !exists {
		logger = newLogger()
		ctx.Set("logger", logger)
		return logger
	}

	logger, ok := l.(*logrus.Logger)
	if !ok {
		return newLogger()
	}

	return logger
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	if env.GetConfig().LogLevel == "debug" {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}
	return logger
}

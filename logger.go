package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	env "github.com/mike-webster/repo-watcher/env"
	"github.com/sirupsen/logrus"
)

func requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func(ctx *gin.Context) logrus.FieldLogger {
			if ctx.Request.URL.Path == "/" && ctx.Writer.Status() == 200 {
				return nil
			}

			// log body if one is given]
			strBody := ""
			body, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				defaultLogger(ctx).WithField("error", err).Error("cant read request body")
			} else {
				// write the body back into the request
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				strBody = string(body)
				strBody = strings.Replace(strBody, "\n", "", -1)
				strBody = strings.Replace(strBody, "\t", "", -1)
			}

			logger := defaultLogger(ctx).WithFields(logrus.Fields{
				"client_ip":    ctx.ClientIP(),
				"event":        "http.in",
				"method":       ctx.Request.Method,
				"path":         ctx.GetString("originalPath"),
				"query":        ctx.Request.URL.RawQuery,
				"referer":      ctx.Request.Referer(),
				"status":       ctx.Writer.Status(),
				"user_agent":   ctx.Request.UserAgent(),
				"git_event":    ctx.Request.Header.Get("X-GitHub-Event"),
				"request_body": strBody,
			})

			if len(ctx.Errors) > 0 {
				logger.Error(strings.TrimSpace(ctx.Errors.String()))
			} else {
				logger.Info()
			}
			return logger
		}(ctx)
	}
}

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

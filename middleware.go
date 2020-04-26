package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// consolidate stack on crahes
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				b, _ := ioutil.ReadAll(c.Request.Body)

				defaultLogger(c).WithFields(logrus.Fields{
					"event":    "ErrPanicked",
					"error":    r,
					"stack":    string(debug.Stack()),
					"path":     c.Request.RequestURI,
					"formbody": string(b),
				}).Error("panic recovered")

				c.AbortWithStatus(500)
			}
		}()
		c.Next() // execute all the handlers
	}
}

func setDependencies(deps *AppDependencies) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("deps", deps)
		ctx.Set("logger", defaultLogger(ctx))
		ctx.Next()
	}
}

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

			ctx.Set("logger", defaultLogger(ctx).WithFields(logrus.Fields{
				"git_event": ctx.Request.Header.Get("X-GitHub-Event"),
			}))

			return logger
		}(ctx)
	}
}

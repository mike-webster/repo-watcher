package main

import (
	"io/ioutil"
	"runtime/debug"

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

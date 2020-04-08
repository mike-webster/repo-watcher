package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const CODE_OK int = 200
const CODE_INVALID int = 400
const CODE_UNAUTH int = 401

func startServer(port int) {
	r := gin.Default()
	r.GET("/", handlerHealtcheck)

	r.Group("/v1")
	{
		r.POST("/github", handlerGitHub)
	}
	r.Run(fmt.Sprintf(":%v", port))
}

func handlerHealtcheck(ctx *gin.Context) {
	ctx.JSON(CODE_OK, fmt.Sprintf("{\"%v\":\"%v\"}", "message", "ok"))
}

func handlerGitHub(ctx *gin.Context) {
	eventName := ctx.Request.Header["X-GitHub-Event"]
	secret := ctx.Request.Header["X-Hub-Signature"]
	body := ctx.Request.PostForm
	message := "Incoming GH WH request:\n\tevent: %v,\n\tsecret: %v, \n\tbody: %v"
	Log(fmt.Sprintf(message, eventName, secret, body), "info")
}

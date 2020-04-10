package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	env "github.com/mike-webster/repo-watcher/env"
	webhookmodels "github.com/mike-webster/repo-watcher/webhookmodels"
)

// CodeOK is for a 200 response
const CodeOK int = 200

// CodeNoContent is for a 204 response
const CodeNoContent int = 204

// CodeInvalid is for a 400 response
const CodeInvalid int = 400

// CodeUnauth is for a 401 response
const CodeUnauth int = 401

const errInvalidSecret = "Invalid request secret"
const errMissingEvent = "Missing event value"
const errInvalidBody = "Invalid POST body"
const errInvalidHeader = "Invalid request headers"

type ApiServer interface {
	Start() error
}

type Server struct {
	Engine *gin.Engine
	Port   string
}

func (s *Server) Run() error {
	return s.Engine.Run(fmt.Sprintf(":%v", s.Port))
}

func SetupServer(port string) *Server {
	router := gin.Default()
	router.GET("/", handlerHealtcheck)

	v1 := router.Group("/v1")
	{
		v1.POST("/github", handlerGitHub)
	}

	return &Server{
		Port:   port,
		Engine: router,
	}
}

func handlerHealtcheck(ctx *gin.Context) {
	ctx.JSON(CodeOK, fmt.Sprintf("{\"%v\":\"%v\"}", "message", "ok"))
}

type ghRequestHeader struct {
	Event  string `header:"X-GitHub-Event" binding:"required"`
	Secret string `header:"X-Hub-Signature" binding:"required"`
}

func (ghrh *ghRequestHeader) ToString() string {
	return fmt.Sprint("Event: ", ghrh.Event, " -- Secret: ", ghrh.Secret)
}

func handlerGitHub(ctx *gin.Context) {
	hdr := &ghRequestHeader{}
	err := ctx.BindHeader(hdr)
	if err != nil {
		Log(fmt.Sprint("invalid request header; error: ", err.Error()), "error")

		errs := strings.Split(err.Error(), "\n")
		msg := ""
		for _, e := range errs {
			v := strings.Replace(e, "Key: ", "", 1)
			msg += fmt.Sprintf("\"%v\":\"%v\",", "reason", v)
		}
		msg = fmt.Sprintf("{%v}", strings.TrimRight(msg, ","))
		ctx.JSON(CodeInvalid, msg)
		return
	}

	summary, errCode, err := parseEventMessage(hdr.Event, ctx)
	if err != nil {
		Log(err.Error(), "error")
		if errCode == CodeInvalid {
			errs := strings.Split(err.Error(), "\n")
			msg := ""
			for _, e := range errs {
				v := strings.Replace(e, "Key: ", "", 1)
				msg += fmt.Sprintf("\"%v\":\"%v\",", "reason", v)
			}
			msg = fmt.Sprintf("{%v}", strings.TrimRight(msg, ","))
			ctx.JSON(CodeInvalid, msg)
			return
		}
	}

	if len(summary) > 0 {
		sendMessageToSlack(summary)
	}
	ctx.Status(CodeNoContent)
}

func parseEventMessage(event string, ctx *gin.Context) (string, int, error) {
	message := ""
	username := ""
	var err error

	if event == "create" {
		event := &webhookmodels.CreateEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "gollum" {
		event := &webhookmodels.GollumEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "issues" {
		event := &webhookmodels.IssuesEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "issue_comment" {
		event := &webhookmodels.IssueCommentEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "project_card" {
		event := &webhookmodels.ProjectCardEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "project_column" {
		event := &webhookmodels.ProjectColumnEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "pull_request" {
		event := &webhookmodels.PullRequestEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "pull_request_review" {
		event := &webhookmodels.PullRequestReviewEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "pull_request_review_comment" {
		event := &webhookmodels.PullRequestReviewCommentEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(err.Error(), "error")
		}
	} else if event == "push" {
		event := &webhookmodels.PushEventPayload{}
		err = ctx.BindJSON(event)
		if err == nil {
			message = event.ToString()
			username = event.Sender.Login
		} else {
			Log(fmt.Sprint("Bad Request: \n", ctx.Request.Form), "error")
		}
	} else {
		Log(fmt.Sprint("Unsupported event: ", event), "error")
	}

	if len(username) > 0 && len(message) > 0 {
		name, err := getNameFromUsername(username)
		if err != nil {
			return fmt.Sprint(name, " ", message), 0, nil
		}
		Log(fmt.Sprint("couldn't parse display name: ", username), "error")
		return fmt.Sprint(username, " ", message), 0, nil
	}

	if err != nil {
		Log(fmt.Sprint("error parsing body for event: ", event, "\n\tError:\t\t", err.Error()), "error")
	}

	return "", CodeInvalid, err
}

func getNameFromUsername(username string) (string, error) {
	cfg := env.GetConfig()
	userBody, err := MakeRequest(fmt.Sprint(cfg.BaseURL(), cfg.UserEndpoint), username, cfg.APIToken)
	if err != nil {
		return "", err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(*userBody, &payload)
	if err != nil {
		return "", err
	}

	name := strings.Split(payload["name"].(string), ", ")
	if len(name) == 2 {
		return fmt.Sprint(name[1], " ", name[0]), nil
	}

	Log(fmt.Sprint("issue finding name: ", username, " -- ", name), "error")

	return username, nil
}

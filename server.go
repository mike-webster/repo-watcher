package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	webhookmodels "github.com/mike-webster/repo-watcher/webhookmodels"
	"github.com/sirupsen/logrus"
)

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

type ghRequestHeader struct {
	Event  string `header:"X-GitHub-Event" binding:"required"`
	Secret string `header:"X-Hub-Signature" binding:"required"`
}

func (ghrh *ghRequestHeader) ToString() string {
	return fmt.Sprint("Event: ", ghrh.Event, " -- Secret: ", ghrh.Secret)
}

// SetupServer will return a configured gin server ready to run on the
// provided port.
func SetupServer(port string, deps *AppDependencies) *Server {
	router := gin.New()
	router.Use(requestLogger())
	router.Use(recovery())
	router.Use(setDependencies(deps))
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

func handlerGitHub(ctx *gin.Context) {
	deps := ctx.MustGet("deps").(*AppDependencies)
	hdr := &ghRequestHeader{}
	err := ctx.BindHeader(hdr)
	if err != nil {
		deps.logger.WithField("error", err).Error("invalid request header -- could not bind")

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

	summary, repo, err := parseEventMessage(ctx, hdr.Event, deps.logger)
	if err != nil {
		deps.logger.WithField("error", err).Error("couldn't parse event message")
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

	if len(summary) > 0 {
		err := deps.dispatchers.ProcessMessage(repo, summary, deps.logger)
		if err != nil {
			deps.logger.WithFields(logrus.Fields{
				"error":   err,
				"payload": summary,
			}).Error("error sending message")
			ctx.Status(500)
			return
		}
	}
	ctx.Status(CodeNoContent)
}

func parseEvent(ctx *gin.Context, eventName string) (webhookmodels.Event, error) {
	switch eventName {
	case "create":
		var event webhookmodels.CreateEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "gollum":
		var event webhookmodels.GollumEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "issue_comment":
		var event webhookmodels.IssueCommentEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "issues":
		var event webhookmodels.IssuesEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "project_card":
		var event webhookmodels.ProjectCardEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "project_column":
		var event webhookmodels.ProjectColumnEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "pull_request":
		var event webhookmodels.PullRequestEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "pull_request_review_comment":
		var event webhookmodels.PullRequestReviewCommentEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "pull_request_review":
		var event webhookmodels.PullRequestReviewEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "push":
		var event webhookmodels.PushEventPayload
		err := ctx.BindJSON(&event)
		if err != nil {
			return nil, err
		}
		return &event, nil
	case "ping":
		// returning an object to force proceesing to be skipped
		return &webhookmodels.CreateEventPayload{
			Repo: webhookmodels.Repository{
				Name: "skip",
			},
		}, nil
	case "status":
		// returning an object to force proceesing to be skipped
		return &webhookmodels.CreateEventPayload{
			Repo: webhookmodels.Repository{
				Name: "skip",
			},
		}, nil
	default:
		defaultLogger(ctx).WithFields(logrus.Fields{
			"event": "unknown_github_event",
			"value": eventName,
		}).Error("unknown event name from github")
		return nil, nil
	}
}

func parseEventMessage(ctx *gin.Context, eventName string, logger *logrus.Logger) (string, string, error) {
	event, err := parseEvent(ctx, eventName)
	if err != nil {
		return "", "", err
	}

	// this is just skipping the initial "ping" for now
	if event.Repository() == "skip" {
		return "", "", nil
	}

	if len(event.ToString()) < 1 {
		logger.WithFields(logrus.Fields{
			"event":      "skipping_notification",
			"event_name": eventName,
		}).Warn("no message returned, skipping notify")
		return "", "", nil
	}

	name, err := getNameFromUsername(event.Username())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event":    "failed_name_retrieval",
			"error":    err,
			"username": event.Username(),
		}).Warn("couldnt retrieve name from username")

		return fmt.Sprint(event.Username(), " ", event.ToString()), event.Repository(), nil
	}

	return fmt.Sprint(name, " ", event.ToString()), event.Repository(), nil
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mike-webster/repo-watcher/env"
	webhookmodels "github.com/mike-webster/repo-watcher/webhookmodels"
	"github.com/sirupsen/logrus"
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

var _logger *logrus.Logger

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
func SetupServer(port string) *Server {
	router := gin.New()
	router.Use(requestLogger())
	router.Use(recovery())
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
	logger := defaultLogger()
	hdr := &ghRequestHeader{}
	err := ctx.BindHeader(hdr)
	if err != nil {
		logger.WithField("error", err).Error("invalid request header -- could not bind")

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

	summary, hook, err := parseEventMessage(ctx, hdr.Event, logger)
	if err != nil {
		logger.WithField("error", err).Error("couldn't parse event message")
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
		sendMessageToSlack(hook, summary)
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
	default:
		defaultLogger().WithFields(logrus.Fields{
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

	cfg := env.GetConfig()
	watcher := cfg.Watchers.Select(event.Repository())
	if watcher == nil {
		logger.WithFields(logrus.Fields{
			"event":        "orphaned_event",
			"github_event": eventName,
			"repository":   event.Repository(),
			"watchers":     cfg.Watchers.ToString(),
		}).Error()
		return "", "", errors.New("orphaned event")
	}

	name, err := getNameFromUsername(event.Username())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event":    "failed_name_retrieval",
			"error":    err,
			"username": event.Username(),
		}).Error("couldnt retrieve name from username")

		return fmt.Sprint(name, " ", event.ToString()), watcher.Webhook, nil
	}

	return fmt.Sprint(event.Username(), " ", event.ToString()), watcher.Webhook, nil
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
				defaultLogger().WithField("error", err).Error("cant read request body")
			} else {
				// write the body back into the request
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				strBody = string(body)
				strBody = strings.Replace(strBody, "\n", "", -1)
				strBody = strings.Replace(strBody, "\t", "", -1)
			}

			logger := defaultLogger().WithFields(logrus.Fields{
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

func defaultLogger() *logrus.Logger {
	if _logger != nil {
		return _logger
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	if env.GetConfig().LogLevel == "debug" {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}
	_logger = logger
	return _logger
}

// consolidate stack on crahes
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				b, _ := ioutil.ReadAll(c.Request.Body)

				defaultLogger().WithFields(logrus.Fields{
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

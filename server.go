package main

import (
	"bytes"
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
		}).Error("couldnt retrieve name from username")

		return fmt.Sprint(event.Username(), " ", event.ToString()), event.Repository(), nil
	}

	return fmt.Sprint(name, " ", event.ToString()), event.Repository(), nil
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
	_logger = logger
	return logger
}

// consolidate stack on crahes
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				b, _ := ioutil.ReadAll(c.Request.Body)

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
		ctx.Next()
	}
}

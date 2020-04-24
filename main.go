package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	dispatchers "github.com/mike-webster/repo-watcher/dispatchers"
	env "github.com/mike-webster/repo-watcher/env"
	models "github.com/mike-webster/repo-watcher/models"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, logger := initApp()

	if cfg.RunType == "solo" {
		deps := AppDependencies{
			dispatchers: getLocalDispatchers(),
			logger:      logger,
		}
		logger.WithField("run_type", "solo").Info()

		sleep := time.Duration(cfg.RefreshTimer) * time.Second
		for {
			runCheck(&deps, logger)
			logger.WithField("sleep_for", sleep).Info()
			time.Sleep(sleep)
		}
	} else if cfg.RunType == "api" {
		logger.WithField("run_type", "api").Info()
		deps := AppDependencies{
			dispatchers: getSlackDispatchers(),
			logger:      logger,
		}

		router := SetupServer(fmt.Sprint(cfg.Port), &deps)
		err := router.Run()
		if err != nil {
			panic(err)
		}
	} else {
		logger.WithField("run_type", "unknown").Error("crashing")
	}
}

func initApp() (*env.Config, *logrus.Logger) {
	cfg := env.GetConfig()
	logger := defaultLogger(nil)
	logger.Info("initializing")

	logger.WithField("watchers", cfg.Watchers).Info()

	return cfg, logger
}

func runCheck(deps *AppDependencies, logger *logrus.Logger) {
	cfg := env.GetConfig()

	logger.WithField("event", "check_history").Debug("checking previous ids")
	ids, err := GetPreviousIDs()
	if err != nil {
		logger.WithField("error", err).Error("could not perform check")
		return
	}

	logger.WithFields(logrus.Fields{
		"event": "id_check",
		"ids":   ids,
	}).Debug("previous ids")

	url := fmt.Sprint(cfg.BaseURL(), fmt.Sprintf(cfg.EventEndpoint, cfg.OrgName, cfg.RepoToWatch))
	eventsBody, err := MakeRequest(url, "", cfg.APIToken)
	if err != nil {
		logger.WithField("error", err).Error("request for events failed")
		return
	}

	var events []models.Event
	err = json.Unmarshal(*eventsBody, &events)
	if err != nil {
		logger.WithField("error", err).Error("couldn't unmarshal events")
		return
	}

	logger.WithField("event", "event_count").Debug(len(events))

	var repoEvents []models.RepositoryEvent
	for _, event := range events {
		repoEvent, err := models.CreateRepositoryEvent(event)
		if err != nil {
			logger.WithField("event", err).Error("coudln't create event from payload")

			continue
		}

		repoEvents = append(repoEvents, repoEvent)
	}

	err = WriteNewIDs(repoEvents)
	if err != nil {
		logger.WithField("error", err).Error("couldn't write event ids")

		// probably shouldn't continue or this may be redundant
		return
	}

	newEvents := []models.RepositoryEvent{}
	for _, e := range repoEvents {
		isOld := false
		for _, o := range *ids {
			if e.Raw().ID == o {
				isOld = true
			}
		}
		if !isOld {
			newEvents = append(newEvents, e)
		}
	}

	if len(newEvents) > 0 {
		// TODO: should we filter out the current user's notifications?
		for _, event := range newEvents {
			logEvent(event.Raw(), logger)
			announceEvent(event, deps, logger)
			time.Sleep(5 * time.Second)
		}
	}
}

func announceEvent(e models.RepositoryEvent, deps *AppDependencies, logger *logrus.Logger) {
	message := e.Say()
	if strings.Contains(message, "#{actor}") {
		realName, err := getNameFromUsername(e.TriggeredBy())
		if err != nil {
			logger.WithField("error", err).Error("couldnt retrieve name from username")

			realName = e.TriggeredBy()
		}

		message = strings.Replace(message, "#{actor}", realName, 1)
		logger.WithField("user_message", message).Info()
	}
	if strings.Contains(message, "#{branch}") {
		message = strings.Replace(message, "#{branch}", e.BranchName(), 1)
	}
	if strings.Contains(message, "#{comment}") {
		message = strings.Replace(message, "#{comment}", e.Comment(), 1)
	}

	// TODO: low priority... I'm not really using solo mode
	// make this able to use more  than just academy
	deps.dispatchers.ProcessMessage("academy", message, deps.logger)
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

	return "", errors.New(fmt.Sprint("could not find name for user: ", username))
}

func logEvent(e models.Event, logger *logrus.Logger) {
	logger.WithFields(logrus.Fields{
		"user":    e.Actor.Username,
		"type":    camelRegexp(e.Type),
		"payload": e.Payload,
	}).Info()
}

func camelRegexp(str string) string {
	re := regexp.MustCompile(`([A-Z]+)`)
	str = re.ReplaceAllString(str, ` $1`)
	str = strings.Trim(str, " ")
	return str
}

func getSlackDispatchers() dispatchers.Dispatchers {
	var ds dispatchers.Dispatchers
	for _, d := range env.GetConfig().Watchers {
		ds = append(ds, &dispatchers.SlackDispatcher{
			URL:      d.Webhook,
			RepoName: d.Repo,
		})
	}
	return ds
}

func getLocalDispatchers() dispatchers.Dispatchers {
	var ds dispatchers.Dispatchers
	for _, d := range env.GetConfig().Watchers {
		ds = append(ds, &dispatchers.LocalDispatcher{
			RepoName: d.Repo,
		})
	}
	return ds
}

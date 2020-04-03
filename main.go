package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	env "github.com/mike-webster/repo-watcher/env"
	models "github.com/mike-webster/repo-watcher/models"
)

func main() {
	cfg := env.GetConfig()
	Log("...starting to monitor...", "info")

	for {
		runCheck()
		time.Sleep(time.Duration(cfg.RefreshTimer*1000000000) * time.Second)
	}

	Log("...stopping monitoring...", "info")
}

func runCheck() {
	cfg := env.GetConfig()

	Log("...checking previous ids...", "debug")
	ids, err := GetPreviousIDs()
	if err != nil {
		panic(err)
	}

	Log(fmt.Sprint("previous ids: ", ids), "info")
	Log("...finding most recents events...", "debug")

	eventsBody, err := MakeRequest(fmt.Sprint(cfg.BaseURL, cfg.EventEndpoint), "academy", cfg.APIToken)
	if err != nil {
		panic(err)
	}

	var events []models.Event
	err = json.Unmarshal(*eventsBody, &events)
	if err != nil {
		panic(err)
	}

	Log(fmt.Sprint("found ", len(events), " events"), "info")

	var repoEvents []models.RepositoryEvent
	for _, event := range events {
		repoEvent, err := models.CreateRepositoryEvent(event)
		if err != nil {
			Log(err.Error(), "error")
		} else {
			repoEvents = append(repoEvents, repoEvent)
		}
	}

	err = WriteNewIDs(repoEvents)
	if err != nil {
		panic(err)
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
			logEvent(event.Raw())
			announceEvent(event)
			time.Sleep(5 * time.Second)
		}
	}
}

func announceEvent(e models.RepositoryEvent) {
	cfg := env.GetConfig()
	message := e.Say()
	if strings.Index(message, "#{user}") > 0 {
		message = strings.Replace(message, "#{user}", cfg.UserName, 1)
	}
	if strings.Index(message, "#{actor}") > 0 {
		realName, err := getNameFromUsername(e.TriggeredBy())
		if err != nil {
			Log(err.Error(), "#error")
			return
		}

		message = strings.Replace(message, "#{actor}", realName, 1)
	}
	if strings.Index(message, "#{branch}") > 0 {
		message = strings.Replace(message, "#{branch}", e.BranchName(), 1)
	}
	if strings.Index(message, "#{comment}") > 0 {
		message = strings.Replace(message, "#{comment}", e.Comment(), 1)
	}

	say(message)
}

func getNameFromUsername(username string) (string, error) {
	cfg := env.GetConfig()
	userBody, err := MakeRequest(fmt.Sprint(cfg.BaseURL, cfg.UserEndpoint), username, cfg.APIToken)
	if err != nil {
		return "", err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(*userBody, &payload)
	if err != nil {
		return "", err
	}

	name := strings.Split(payload["name"].(string), ", ")
	return fmt.Sprint(name[1], " ", name[0]), nil
}

func logEvent(e models.Event) {

	Log(fmt.Sprint("User ", e.Actor.Username), "debug")
	Log(fmt.Sprint("Event type: ", camelRegexp(e.Type)), "info")
	Log(fmt.Sprint("Payload: ", e.Payload), "info")
}

func camelRegexp(str string) string {
	re := regexp.MustCompile(`([A-Z]+)`)
	str = re.ReplaceAllString(str, ` $1`)
	str = strings.Trim(str, " ")
	return str
}

func say(message string) {
	cmd := exec.Command("say", message)
	cmd.Run()
}

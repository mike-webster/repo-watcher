package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	env "github.com/mike-webster/repo-watcher/env"
	models "github.com/mike-webster/repo-watcher/models"
)

func main() {
	Log("...starting to monitor...", "info")
	runCheck()
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

	url := fmt.Sprint(cfg.BaseURL(), fmt.Sprintf(cfg.EventEndpoint, cfg.OrgName, cfg.RepoToWatch))
	eventsBody, err := MakeRequest(url, "", cfg.APIToken)
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
	message := e.Say()
	if strings.Contains(message, "#{actor}") {
		realName, err := getNameFromUsername(e.TriggeredBy())
		if err != nil {
			Log(err.Error(), "#error")
			return
		}

		message = strings.Replace(message, "#{actor}", realName, 1)
		Log(fmt.Sprint("User message:", message), "info")
	}
	if strings.Contains(message, "#{branch}") {
		message = strings.Replace(message, "#{branch}", e.BranchName(), 1)
	}
	if strings.Contains(message, "#{comment}") {
		message = strings.Replace(message, "#{comment}", e.Comment(), 1)
	}

	sendMessageToSlack(message)
}

func sendMessageToSlack(message string) error {
	cfg := env.GetConfig()
	body := fmt.Sprintf("{\"text\":\"%v\"}", message)
	req, err := http.NewRequest("POST", cfg.SlackWebhook, strings.NewReader(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprint("non-200 response: ", resp.StatusCode))
	}
	return nil
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

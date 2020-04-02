package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	env "github.com/mike-webster/repo-watcher/env"
	models "github.com/mike-webster/repo-watcher/models"
)

func main() {
	Log("...starting to monitor...", "info")

	for {
		runCheck()
		time.Sleep(5 * time.Minute)
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

	err = WriteNewIDs(events)
	if err != nil {
		panic(err)
	}

	newEvents := []models.Event{}
	for _, e := range events {
		isOld := false
		for _, o := range *ids {
			if e.ID == o {
				Log(fmt.Sprintf("%v == %v", e.ID, o), "debug")
				isOld = true
			}
		}
		if !isOld {
			newEvents = append(newEvents, e)
		}
	}

	if len(newEvents) > 0 {
		cmd := exec.Command("say", fmt.Sprintf("Hey, %v ! We just found %v new notifications in your %v repo.", cfg.UserName, len(newEvents), "Academy"))
		cmd.Run()
	}
}

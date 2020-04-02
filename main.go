package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	Log("...starting to monitor...", "info")
	cfg := env.GetConfig()

	eventsBody, err := MakeRequest(fmt.Sprint(cfg.BaseURL, cfg.EventEndpoint), "academy", cfg.APIToken)
	if err != nil {
		panic(err)
	}

	var events []models.Event
	err = json.Unmarshal(eventsBody, &events)
	if err != nil {
		panic(err)
	}

	Log("...stopping monitoring", "info")
}

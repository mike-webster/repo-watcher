package webhookmodels

import (
	"fmt"
	"strings"
)

// GollumEventPayload is the request receiced when a wiki page is created
// or updated.
//
// https://developer.github.com/v3/activity/events/types/#gollumevent
type GollumEventPayload struct {
	Pages  []Page     `json:"pages"  binding:"required"`
	Repo   Repository `json:"repository"`
	Sender User       `json:"sender"`
}

// Names will return the names of each page edited
func (gep *GollumEventPayload) Names() string {
	ret := []string{}
	for _, p := range gep.Pages {
		ret = append(ret, p.Name)
	}
	return strings.Join(ret, "\n")
}

// ToString outputs a summary message of the event
func (gep *GollumEventPayload) ToString() string {
	return fmt.Sprintf("updated some wiki content: \n%v", gep.Names())
}

// Username returns the username of the user who triggered the event
func (gep *GollumEventPayload) Username() string {
	return gep.Sender.Login
}

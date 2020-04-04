package models

import (
	"fmt"
	"strings"
)

// GollumEvent is triggered when a Wiki page is created or updated.
type GollumEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (ge *GollumEvent) TriggeredBy() string {
	return ge.raw.Actor.Username
}

// EventType returns the type of event
func (ge *GollumEvent) EventType() string {
	return "GollumEvent"
}

// BranchName returns the name of the branch, if there is one
func (ge *GollumEvent) BranchName() string {
	return ge.raw.Payload["ref"].(string)
}

// Comment returns the comment, if there is one
func (ge *GollumEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (ge *GollumEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (ge *GollumEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (ge *GollumEvent) Path() string {
	return ""
}

// Raw returns the underlying Event ojbect
func (ge *GollumEvent) Raw() Event {
	return ge.raw
}

// Say returns the templated string to pass to the say command for this object
func (ge *GollumEvent) Say() string {
	pagesBody := ge.Raw().Payload["pages"].([]interface{})
	pages := []string{}
	for _, page := range pagesBody {
		iPage := page.(map[string]interface{})
		pages = append(pages, iPage["html_url"].(string))
	}

	return fmt.Sprint("#{actor} just edited some Wiki content:\n", strings.Join(pages, "\n"))
}

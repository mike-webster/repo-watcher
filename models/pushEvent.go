package models

import "fmt"

// PushEvent represents a user pushing code to a repository
type PushEvent struct {
	raw Event
}

// TriggeredBy returns the username of the user who triggered the event
func (pe *PushEvent) TriggeredBy() string {
	return pe.raw.Actor.Username
}

// EventType returns the type of event
func (pe *PushEvent) EventType() string {
	return "PushEvent"
}

// BranchName returns the name of the branch, if there is one
func (pe *PushEvent) BranchName() string {
	longName := pe.raw.Payload["ref"].(string)
	return parseBranchName(longName)
}

// Comment returns the comment, if there is one
func (pe *PushEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (pe *PushEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (pe *PushEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (pe *PushEvent) Path() string {
	return pe.raw.Payload["ref"].(string)
}

// Raw returns the underlying Event ojbect
func (pe *PushEvent) Raw() Event {
	return pe.raw
}

// Say returns the templated string to pass to the say command for this object
func (pe *PushEvent) Say() string {
	return fmt.Sprint("Hey, #{user}! #{actor} just pushed some code to branch #{branch}")
}

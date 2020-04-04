package models

import "fmt"

// IssuesEvent is triggered by an issue being opened, edited, deleted, pinned,
// unpinned, closed, reopened, assigned, unassigned, labeled, unlabeled,
// locked, unlocked, transferred, milestoned, or demilestoned.
type IssuesEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (ie *IssuesEvent) TriggeredBy() string {
	return ie.raw.Actor.Username
}

// EventType returns the type of event
func (ie *IssuesEvent) EventType() string {
	return "IssuesEvent"
}

// BranchName returns the name of the branch, if there is one
func (ie *IssuesEvent) BranchName() string {
	return ""
}

// Comment returns the comment, if there is one
func (ie *IssuesEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (ie *IssuesEvent) Title() string {
	issue := ie.raw.Payload["issue"].(map[string]interface{})
	return issue["title"].(string)
}

// Body returns the body, if there is one
func (ie *IssuesEvent) Body() string {
	issue := ie.raw.Payload["issue"].(map[string]interface{})
	return issue["body"].(string)
}

// Path returns the html path to the content, if there is one
func (ie *IssuesEvent) Path() string {
	issue := ie.raw.Payload["issue"].(map[string]interface{})
	return issue["html_url"].(string)
}

// Raw returns the underlying Event ojbect
func (ie *IssuesEvent) Raw() Event {
	return ie.raw
}

// Say returns the templated string to pass to the say command for this object
func (ie *IssuesEvent) Say() string {
	action := ie.Raw().Payload["action"].(string)
	issue := ie.Raw().Payload["issue"].(map[string]interface{})
	title := issue["title"]
	return fmt.Sprint("#{actor} just ", action, " an issue with the title: ", title)
}

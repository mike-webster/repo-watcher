package models

import "fmt"

// IssueCommentEvent is triggered when an issue comment is created, edited, or deleted.
type IssueCommentEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (ice *IssueCommentEvent) TriggeredBy() string {
	return ice.raw.Actor.Username
}

// EventType returns the type of event
func (ice *IssueCommentEvent) EventType() string {
	return "IssueCommentEvent"
}

// BranchName returns the name of the branch, if there is one
func (ice *IssueCommentEvent) BranchName() string {
	return ""
}

// Comment returns the comment, if there is one
func (ice *IssueCommentEvent) Comment() string {
	comment := ice.raw.Payload["comment"].(map[string]interface{})
	return comment["body"].(string)
}

// Title returns the title, if there is one
func (ice *IssueCommentEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (ice *IssueCommentEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (ice *IssueCommentEvent) Path() string {
	comment := ice.raw.Payload["comment"].(map[string]interface{})
	return comment["url"].(string)
}

// Raw returns the underlying Event ojbect
func (ice *IssueCommentEvent) Raw() Event {
	return ice.raw
}

// Say returns the templated string to pass to the say command for this object
func (ice *IssueCommentEvent) Say() string {
	issue := ice.Raw().Payload["issue"].(map[string]interface{})
	title := issue["title"]
	return fmt.Sprint("#{actor} just commented on an issue with the title: ", title, ", Here's the comment: #{comment}")
}

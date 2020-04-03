package models

import "fmt"

// PullRequestEvent represents a user acting on a pull request
type PullRequestEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (pre *PullRequestEvent) TriggeredBy() string {
	return pre.raw.Actor.Username
}

// EventType returns the type of event
func (pre *PullRequestEvent) EventType() string {
	return "PullRequestEvent"
}

// BranchName returns the name of the branch, if there is one
func (pre *PullRequestEvent) BranchName() string {
	pr := pre.Raw().Payload["pull_request"].(map[string]interface{})
	head := pr["head"].(map[string]interface{})
	return head["ref"].(string)
}

// Comment returns the comment, if there is one
func (pre *PullRequestEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (pre *PullRequestEvent) Title() string {
	pr := pre.Raw().Payload["pull_request"].(map[string]interface{})
	return pr["title"].(string)
}

// Body returns the body, if there is one
func (pre *PullRequestEvent) Body() string {
	pr := pre.Raw().Payload["pull_request"].(map[string]interface{})
	return pr["body"].(string)
}

// Path returns the html path to the content, if there is one
func (pre *PullRequestEvent) Path() string {
	pr := pre.Raw().Payload["pull_request"].(map[string]interface{})
	return pr["html_url"].(string)
}

// Raw returns the underlying Event ojbect
func (pre *PullRequestEvent) Raw() Event {
	return pre.raw
}

// Say returns the templated string to pass to the say command for this object
func (pre *PullRequestEvent) Say() string {
	action := pre.Raw().Payload["action"].(string)
	pr := pre.Raw().Payload["pull_request"].(map[string]interface{})
	title := pr["title"]
	return fmt.Sprint("Hey, #{user}! #{actor} just ", action, " a pull request with the title: ", title)
}

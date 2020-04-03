package models

import "fmt"

// PullRequestReviewCommentEvent represents an actor creating something
type PullRequestReviewCommentEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (prrce *PullRequestReviewCommentEvent) TriggeredBy() string {
	return prrce.raw.Actor.Username
}

// EventType returns the type of event
func (prrce *PullRequestReviewCommentEvent) EventType() string {
	return "PullRequestReviewCommentEvent"
}

// BranchName returns the name of the branch, if there is one
func (prrce *PullRequestReviewCommentEvent) BranchName() string {
	return ""
}

// Comment returns the comment, if there is one
func (prrce *PullRequestReviewCommentEvent) Comment() string {
	comment := prrce.raw.Payload["comment"].(map[string]interface{})
	return comment["body"].(string)
}

// Title returns the title, if there is one
func (prrce *PullRequestReviewCommentEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (prrce *PullRequestReviewCommentEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (prrce *PullRequestReviewCommentEvent) Path() string {
	comment := prrce.raw.Payload["comment"].(map[string]interface{})
	return comment["html_url"].(string)
}

// Raw returns the underlying Event ojbect
func (prrce *PullRequestReviewCommentEvent) Raw() Event {
	return prrce.raw
}

// Say returns the templated string to pass to the say command for this object
func (prrce *PullRequestReviewCommentEvent) Say() string {
	comment := prrce.Raw().Payload["comment"].(map[string]interface{})
	file := comment["path"].(string)
	pr := prrce.Raw().Payload["pull_request"].(map[string]interface{})
	prTitle := pr["title"]
	return fmt.Sprint("Hey, #{user}! #{actor} just commented on a pull request with the title: ", prTitle, ", The comment was in the file: ", file, ", Here's the comment: #{comment}")
}

package webhookmodels

import "fmt"

// IssueCommentEventPayload is the request received when an issue comment is
// created, edited, or deleted.
//
// https://developer.github.com/v3/activity/events/types/#issuecommentevent
type IssueCommentEventPayload struct {
	Action  string      `json:"action"  binding:"required"`
	Issue   Issue       `json:"issue"`
	Comment interface{} `json:"comment"`
	Repo    Repository  `json:"repository"`
	Sender  User        `json:"sender"`
}

// ToString outputs a summary message of the event
func (icep *IssueCommentEventPayload) ToString() string {
	return fmt.Sprintf("%v on an issue: \n----\n| Title: %v\n----\nComment: \n%v", icep.Action, icep.Issue.Title, icep.Comment)
}

// Username returns the username of the user who triggered the event
func (icep *IssueCommentEventPayload) Username() string {
	return icep.Sender.Login
}

func (icep *IssueCommentEventPayload) Repository() string {
	return icep.Repo.Name
}

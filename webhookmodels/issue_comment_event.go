package webhookmodels

import "fmt"

import "github.com/mike-webster/repo-watcher/markdown"

// IssueCommentEventPayload is the request received when an issue comment is
// created, edited, or deleted.
//
// https://developer.github.com/v3/activity/events/types/#issuecommentevent
type IssueCommentEventPayload struct {
	Action  string       `json:"action"  binding:"required"`
	Issue   Issue        `json:"issue"`
	Comment IssueComment `json:"comment"`
	Repo    Repository   `json:"repository"`
	Sender  User         `json:"sender"`
}

type IssueComment struct {
	Body string `json:"body"`
}

// ToString outputs a summary message of the event
func (icep *IssueCommentEventPayload) ToString() string {
	header := markdown.MarkdownBold(fmt.Sprintf("%s on an issue", icep.Action))
	title := markdown.MarkdownLink(icep.Issue.URL, fmt.Sprintf("Title: %s", icep.Issue.Title))
	comment := markdown.MarkdownMultilineCode(fmt.Sprintf(icep.Comment.Body))
	return fmt.Sprintf("%s\n%s\n%s", header, title, comment)
}

// Username returns the username of the user who triggered the event
func (icep *IssueCommentEventPayload) Username() string {
	return icep.Sender.Login
}

func (icep *IssueCommentEventPayload) Repository() string {
	return icep.Repo.Name
}

package webhookmodels

import "fmt"

import "github.com/mike-webster/repo-watcher/markdown"

// PullRequestReviewCommentEventPayload is the request received when a comment
// on a pull request's unified dif is created, edited, or deleted.
//
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
type PullRequestReviewCommentEventPayload struct {
	Action      string        `json:"action" binding:"required"`
	PullRequest PullRequest   `json:"pull_request"`
	Comment     ReviewComment `json:"comment"`
	Repo        Repository    `json:"repository"`
	Sender      User          `json:"sender"`
}

// ToString outputs a summary message of the event
func (prrcep *PullRequestReviewCommentEventPayload) ToString() string {
	header := markdown.MarkdownBold(fmt.Sprintf("%s a comment on a pull request review", prrcep.Action))
	title := markdown.MarkdownLink(prrcep.PullRequest.URL, fmt.Sprintf("Title:  %s", prrcep.PullRequest.Title))
	comment := markdown.MarkdownMultilineCode(prrcep.Comment.Body)
	return fmt.Sprintf("%s\n%s\n%s", header, title, comment)
}

// Username returns the username of the user who triggered the event
func (prrcep *PullRequestReviewCommentEventPayload) Username() string {
	return prrcep.Sender.Login
}

func (prrcep *PullRequestReviewCommentEventPayload) Repository() string {
	return prrcep.Repo.Name
}

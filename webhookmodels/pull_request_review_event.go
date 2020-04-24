package webhookmodels

import "fmt"

import "github.com/mike-webster/repo-watcher/markdown"

// PullRequestReviewEventPayload is the request received when a pull request review
// is submitted into a non-pending state, the body is edited, or the review is dismissed.
//
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewevent
type PullRequestReviewEventPayload struct {
	Action      string      `json:"action" binding:"required"`
	PullRequest PullRequest `json:"pull_request"`
	Review      Review      `json:"review"`
	Repo        Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

// ToString outputs a summary message of the event
func (prrep *PullRequestReviewEventPayload) ToString() string {
	header := markdown.MarkdownBold(fmt.Sprintf("%s a pull request review", prrep.Action))
	title := markdown.MarkdownLink(prrep.PullRequest.URL, fmt.Sprintf("Title: %s", prrep.PullRequest.Title))
	state := fmt.Sprintf("State: %s", prrep.Review.State)
	body := markdown.MarkdownMultilineCode(fmt.Sprintf("Body: \n%v", prrep.PullRequest.Body))
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, title, state, body)
}

// Username returns the username of the user who triggered the event
func (prrep *PullRequestReviewEventPayload) Username() string {
	return prrep.Sender.Login
}

func (prrep *PullRequestReviewEventPayload) Repository() string {
	return prrep.Repo.Name
}

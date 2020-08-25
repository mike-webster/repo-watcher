package webhookmodels

import (
	"fmt"

	"github.com/mike-webster/repo-watcher/markdown"
)

// PullRequestEventPayload is the request received when a pull request is assigned,
// unassigned, labeled, unlabeled, opened, edited, closed, reopened, synchronize,
// ready_for_review, locked, unlocked or when a pull review is requested or removed.
//
// https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEventPayload struct {
	Action      string      `json:"action"  binding:"required"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repo        Repository  `json:"repository"`
	Sender      User        `json:"sender"`
	Merged      bool        `json:"merged"`
	Additions   int         `json:"additions"`
	Deletions   int         `json:"deletions"`
}

// ToString outputs a summary message of the event
func (prep *PullRequestEventPayload) ToString() string {
	header := markdown.MarkdownBold(fmt.Sprintf("%v a pull request", prep.Action))
	title := markdown.MarkdownLink(prep.PullRequest.URL, fmt.Sprintf("Title: %s", prep.PullRequest.Title))
	body := markdown.MarkdownMultilineCode(prep.PullRequest.Body)
	if prep.Action == "opened" || prep.Action == "edited" {
		return fmt.Sprintf("%s\n%s\n%s", header, title, body)
	} else if prep.Action == "labeled" {
		labels := markdown.MarkdownMultilineCode(markdown.MarkdownList(prep.PullRequest.Labels.Names()))
		return fmt.Sprintf("%s\n%s\nLabels:\n%s", header, title, labels)
	} else if prep.Action == "closed" {
		return fmt.Sprintf("%s\n%s", header, title)
	}

	return ""
}

// Username returns the username of the user who triggered the event
func (prep *PullRequestEventPayload) Username() string {
	return prep.Sender.Login
}

func (prep *PullRequestEventPayload) Repository() string {
	return prep.Repo.Name
}

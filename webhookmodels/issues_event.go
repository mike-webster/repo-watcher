package webhookmodels

import "fmt"

import "github.com/mike-webster/repo-watcher/markdown"

// IssuesEventPayload is the request received when an issue is opened, edited,
// deleted, pinned, unpinned, closed, reopened, assigned, unassigned, labeled,
// unlabeled, locked, unlocked, transferred, milestoned, or demilestoned
//
// https://developer.github.com/v3/activity/events/types/#issuesevent
type IssuesEventPayload struct {
	Action  string      `json:"action"  binding:"required"`
	Issue   Issue       `json:"issue"`
	Changes interface{} `json:"changes"`
	Repo    Repository  `json:"repository"`
	Sender  User        `json:"sender"`
}

// ToString outputs a summary message of the event
func (iep *IssuesEventPayload) ToString() string {
	header := markdown.MarkdownBold(fmt.Sprintf("%v an issue", iep.Action))
	title := markdown.MarkdownLink(iep.Issue.URL, fmt.Sprintf("Title: %v", iep.Issue.Title))
	body := markdown.MarkdownMultilineCode(iep.Issue.Body)
	return fmt.Sprintf("%s\n%s\n%s", header, title, body)
}

// Username returns the username of the user who triggered the event
func (iep *IssuesEventPayload) Username() string {
	return iep.Sender.Login
}

func (iep *IssuesEventPayload) Repository() string {
	return iep.Repo.Name
}

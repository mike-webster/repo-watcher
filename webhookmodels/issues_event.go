package webhookmodels

import "fmt"

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
	return fmt.Sprintf("%v an issue: \n----\n| Title: %v\n----\nBody: \n%v", iep.Action, iep.Issue.Title, iep.Issue.Body)
}

// Username returns the username of the user who triggered the event
func (iep *IssuesEventPayload) Username() string {
	return iep.Sender.Login
}

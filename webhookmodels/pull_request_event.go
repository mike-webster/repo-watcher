package webhookmodels

import "fmt"

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
}

// ToString outputs a summary message of the event
func (prep *PullRequestEventPayload) ToString() string {
	return fmt.Sprintf("%v a pull request: \n----\n| Title: %v\n----\n| Body: \n%v", prep.Action, prep.PullRequest.Title, prep.PullRequest.Body)
}

// Username returns the username of the user who triggered the event
func (prep *PullRequestEventPayload) Username() string {
	return prep.Sender.Login
}

func (prep *PullRequestEventPayload) Repository() string {
	return prep.Repo.Name
}

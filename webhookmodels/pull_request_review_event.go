package webhookmodels

import "fmt"

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
	return fmt.Sprintf("%v a pull request review for %v\n\t\tNew review state: %v", prrep.Action, prrep.PullRequest.Title, prrep.PullRequest.State)
}

// Username returns the username of the user who triggered the event
func (prrep *PullRequestReviewEventPayload) Username() string {
	return prrep.Sender.Login
}

package webhookmodels

import "fmt"

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
	return fmt.Sprintf("%v a comment on a pull request review for %v\n\t\tNew review state: %v\nComment: %v", prrcep.Action, prrcep.PullRequest.Title, prrcep.PullRequest.State, prrcep.Comment.Body)
}

// Username returns the username of the user who triggered the event
func (prrcep *PullRequestReviewCommentEventPayload) Username() string {
	return prrcep.Sender.Login
}

func (prrcep *PullRequestReviewCommentEventPayload) Repository() string {
	return prrcep.Repo.Name
}

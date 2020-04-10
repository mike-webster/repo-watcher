package webhookmodels

import "time"

// ReviewComment represents a user's comment on a pull request
type ReviewComment struct {
	ID             int64     `json:"id"`
	ReviewID       int64     `json:"pull_request_review_id"`
	NodeID         string    `json:"node_id"`
	Path           string    `json:"path"`
	URL            string    `json:"html_url"`
	User           User      `json:"user"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PullRequestURL string    `json:"pull_request_url"`
}

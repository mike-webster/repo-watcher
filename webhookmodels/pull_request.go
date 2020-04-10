package webhookmodels

import "time"

// PullRequest represents a github pull request
type PullRequest struct {
	ID        int64      `json:"id"`
	NodeID    string     `json:"node_id"`
	URL       string     `json:"html_url"`
	Number    int        `json:"number"`
	State     string     `json:"state"`
	Title     string     `json:"title"`
	User      User       `json:"user"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Assignee  User       `json:"assignee"`
	Labels    []Label    `json:"labels"`
	Repo      Repository `json:"repository"`
	Sender    User       `json:"sender"`
}

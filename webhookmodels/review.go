package webhookmodels

import "time"

// Review repesents a github pull request review
type Review struct {
	ID          int64     `json:"id"`
	NodeID      string    `json:"node_id"`
	User        User      `json:"user"`
	Body        string    `json:"body"`
	SubmittedAt time.Time `json:"submitted_at"`
	State       string    `json:"state"`
	URL         string    `json:"html_url"`
}

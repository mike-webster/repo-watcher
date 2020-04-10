package webhookmodels

import (
	"time"
)

// Comment represets an issue comment
type Comment struct {
	ID        int64     `json:"id"`
	URL       string    `json:"html_url"`
	NodeID    string    `json:"node_id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

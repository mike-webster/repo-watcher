package webhookmodels

import (
	"time"
)

// Card represents a card on a project board
type Card struct {
	ID        int64     `json:"id"`
	URL       string    `json:"html_url"`
	ColumnID  int64     `json:"column_id"`
	Note      string    `json:"note"`
	Creator   User      `json:"creator"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

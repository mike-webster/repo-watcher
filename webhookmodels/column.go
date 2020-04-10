package webhookmodels

// Column represents a column from a project board
type Column struct {
	ID     int64      `json:"id"`
	NodeID string     `json:"node_id"`
	Name   string     `json:"name"`
	Repo   Repository `json:"repository"`
	Sender User       `json:"sender"`
}

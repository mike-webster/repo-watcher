package webhookmodels

// Repository represents a github repository
type Repository struct {
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	Name        string `json:"name"`
	Owner       User   `json:"owner"`
	Sender      User   `json:"sender"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
}

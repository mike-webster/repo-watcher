package webhookmodels

// User represents a github user
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"html_url"`
}

package webhookmodels

// Issue represents a github issue
type Issue struct {
	ID       int64   `json:"id"`
	URL      string  `json:"html_url"`
	NodeID   string  `json:"node_id"`
	Title    string  `json:"title"`
	User     User    `json:"user"`
	Labels   []Label `json:"labels"`
	State    string  `json:"state"`
	Assignee User    `json:"assignee"`
	Body     string  `json:"body"`
}

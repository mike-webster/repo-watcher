package webhookmodels

// Label represents a github label
type Label struct {
	ID     int64  `json:"id"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

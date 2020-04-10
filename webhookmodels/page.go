package webhookmodels

// Page represents a wiki page in a repo
type Page struct {
	Name    string `json:"page_name"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Action  string `json:"action"`
	SHA     string `json:"sha"`
	URL     string `json:"html_url"`
}

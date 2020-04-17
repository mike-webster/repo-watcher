package webhookmodels

// Event represents a payload from the github repo
type Event interface {
	ToString() string
	Username() string
}

package models

// CreateEvent represents an actor creating something
type CreateEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (ce *CreateEvent) TriggeredBy() string {
	return ce.raw.Actor.Username
}

// EventType returns the type of event
func (ce *CreateEvent) EventType() string {
	return "CreateEvent"
}

// BranchName returns the name of the branch, if there is one
func (ce *CreateEvent) BranchName() string {
	return ce.raw.Payload["ref"].(string)
}

// Comment returns the comment, if there is one
func (ce *CreateEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (ce *CreateEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (ce *CreateEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (ce *CreateEvent) Path() string {
	return ""
}

// Raw returns the underlying Event ojbect
func (ce *CreateEvent) Raw() Event {
	return ce.raw
}

// Say returns the templated string to pass to the say command for this object
func (ce *CreateEvent) Say() string {
	return "Hey, #{user}! #{actor} just created something. Branch name: #{branch}"
}

package models

// NOTES: I don't see this coming up. Maybe the project board notifications
//        aren't coming through the way I'm requesting noticications?

// ProjectCardEvent is triggered when a project card is created, edited,
// moved, converted to an issue, or deleted.
type ProjectCardEvent struct {
	raw Event
}

// TriggeredBy returns the name of the username who triggered the event
func (pce *ProjectCardEvent) TriggeredBy() string {
	return pce.raw.Actor.Username
}

// EventType returns the type of event
func (pce *ProjectCardEvent) EventType() string {
	return "ProjectCardEvent"
}

// BranchName returns the name of the branch, if there is one
func (pce *ProjectCardEvent) BranchName() string {
	return ""
}

// Comment returns the comment, if there is one
func (pce *ProjectCardEvent) Comment() string {
	return ""
}

// Title returns the title, if there is one
func (pce *ProjectCardEvent) Title() string {
	return ""
}

// Body returns the body, if there is one
func (pce *ProjectCardEvent) Body() string {
	return ""
}

// Path returns the html path to the content, if there is one
func (pce *ProjectCardEvent) Path() string {
	return ""
}

// Raw returns the underlying Event ojbect
func (pce *ProjectCardEvent) Raw() Event {
	return pce.raw
}

// Say returns the templated string to pass to the say command for this object
func (pce *ProjectCardEvent) Say() string {
	return "TODO"
}

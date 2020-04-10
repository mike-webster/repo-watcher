package webhookmodels

import "fmt"

// ProjectColumnEventPayload is the request received when a project column
// is created, updated, moved, or deleted.
//
// https://developer.github.com/v3/activity/events/types/#projectcolumnevent
type ProjectColumnEventPayload struct {
	Action string     `json:"action"  binding:"required"`
	Column Column     `json:"project_column"`
	Repo   Repository `json:"repository"`
	Sender User       `json:"sender"`
}

// ToString outputs a summary message of the event
func (pcep *ProjectColumnEventPayload) ToString() string {
	return fmt.Sprintf("%v a column: \n----\nName: %v", pcep.Action, pcep.Column.Name)
}

package webhookmodels

import "fmt"

// ProjectCardEventPayload is the request received when a project card is created,
// edited, moved, converted to an issue, or deleted.
//
// https://developer.github.com/v3/activity/events/types/#projectcardevent
type ProjectCardEventPayload struct {
	Action string     `json:"action"  binding:"required"`
	Card   Card       `json:"project_card"`
	Repo   Repository `json:"repository"`
	Sender User       `json:"sender"`
}

// ToString outputs a summary message of the event
func (pcep *ProjectCardEventPayload) ToString() string {
	return fmt.Sprintf("%v a card: \n----\nNote: %v", pcep.Action, pcep.Card.Note)
}

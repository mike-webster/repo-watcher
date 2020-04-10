package webhookmodels

import (
	"fmt"
)

// CreateEventPayload is the request received when a branch or tag is created
// in a repository.
//
// https://developer.github.com/v3/activity/events/types/#createevent
type CreateEventPayload struct {
	Type        string     `json:"ref_type"`
	Ref         string     `json:"ref"  binding:"required"`
	Description string     `json:"description"`
	Repo        Repository `json:"repository"`
	Sender      User       `json:"sender"`
}

// ToString outputs a summary message of the event
func (cep *CreateEventPayload) ToString() string {
	return fmt.Sprintf("created a branch: %v", cep.Ref)
}

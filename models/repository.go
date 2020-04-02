package models

import (
	"fmt"
	"time"
)

// Repository contains information about a GitHub repository
type Repository struct {
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
}

// Actor contains information about a user who performed an event
// in a GitHub repository.
type Actor struct {
	ID        int    `json:"id"`
	Username  string `json:"display_login"`
	AvatarURL string `json:"avatar_url"`
}

// Event contains information about a RepositoryEvent GitHub object
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Actor     Actor                  `json:"actor"`
	Repo      Repository             `json:"repo"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}

// TrackingStub returns a json tag to use to track the events
func (e *Event) TrackingStub() string {
	return fmt.Sprintf("{'id':'%v'}", e.ID)
}

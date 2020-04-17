package webhookmodels

import (
	"fmt"
	"strings"
)

// PushEventPayload is the request received there is a push to a repository branch.
//
// https://developer.github.com/v3/activity/events/types/#pushevent
type PushEventPayload struct {
	Ref     string        `json:"ref" binding:"required"`
	Commits []interface{} `json:"commits"`
	Repo    Repository    `json:"repository"`
	URL     string        `json:"html_url"`
	Sender  User          `json:"sender"`
}

// CommitMessages returns a formatted list of strings from all commits
// in the push
func (pep *PushEventPayload) CommitMessages() string {
	ret := []string{}
	for _, c := range pep.Commits {
		body := c.(map[string]interface{})
		message := body["message"].(string)
		ret = append(ret, message)
	}
	return strings.Join(ret, "\n")
}

// ToString outputs a summary message of the event
func (pep *PushEventPayload) ToString() string {
	return fmt.Sprintf("pushed some changes to %v\nCommits:\n> %v", pep.Ref, pep.CommitMessages())
}

// Username returns the username of the user who triggered the event
func (pep *PushEventPayload) Username() string {
	return pep.Sender.Login
}

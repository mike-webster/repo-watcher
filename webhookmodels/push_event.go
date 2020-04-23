package webhookmodels

import (
	"errors"
	"fmt"
	"github.com/mike-webster/repo-watcher/markdown"
	"reflect"
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
func (pep *PushEventPayload) CommitMessages() (string, *[]error) {
	ret := []string{}
	errs := []error{}
	for _, c := range pep.Commits {
		body, ok := c.(map[string]interface{})
		if !ok {
			errs = append(errs, errors.New("couldn't parse commits"))
			continue
		}

		message, ok := body["message"].(string)
		if !ok {
			errs = append(errs, errors.New(fmt.Sprint("bad conversion: ", reflect.TypeOf(c))))
			continue
		}
		ret = append(ret, fmt.Sprint(message))
	}
	if len(errs) < 1 {
		return strings.Join(ret, "\n"), nil
	}

	return strings.Join(ret, "\n"), &errs
}

// ToString outputs a summary message of the event
func (pep *PushEventPayload) ToString() string {
	messages, errs := pep.CommitMessages()
	if errs != nil {
		for _, m := range *errs {
			fmt.Println("commit message parse error: ", m.Error())
		}
	}
	header := markdown.MarkdownLink(pep.URL, fmt.Sprintf("pushed some changes to %s", pep.Ref))
	title := markdown.MarkdownItalic("Commits")
	body := markdown.MarkdownCode(messages)
	return fmt.Sprintf("%s\n%s\n%s", header, title, body)
}

// Username returns the username of the user who triggered the event
func (pep *PushEventPayload) Username() string {
	return pep.Sender.Login
}

func (pep *PushEventPayload) Repository() string {
	return pep.Repo.Name
}

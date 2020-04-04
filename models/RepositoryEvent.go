package models

import (
	"errors"
	"fmt"
)

// https://developer.github.com/v3/activity/events/types

// RepositoryEvent is a representation of an Event object from the
// GitHub notifications API.
type RepositoryEvent interface {
	// TriggeredBy returns the username of the user who triggered the event
	TriggeredBy() string
	// EventType returns the type of event
	EventType() string
	// BranchName returns the name of the branch, if there is one
	BranchName() string
	// Comment returns the comment, if there is one
	Comment() string
	// Title returns the title, if there is one
	Title() string
	// Body returns the body, if there is one
	Body() string
	// Path returns the html path to the content, if there is one
	Path() string
	// Raw returns the underlying Event ojbect
	Raw() Event
	// Say returns the templated string to pass to the say command for this object
	Say() string
}

// CreateRepositoryEvent will take an event and wrap it with the appropriate
// struct.
func CreateRepositoryEvent(e Event) (RepositoryEvent, error) {
	switch e.Type {
	case "CreateEvent":
		return &CreateEvent{raw: e}, nil
	case "IssueCommentEvent":
		return &IssueCommentEvent{raw: e}, nil
	case "IssuesEvent":
		return &IssuesEvent{raw: e}, nil
	case "PullRequestReviewCommentEvent":
		return &PullRequestReviewCommentEvent{raw: e}, nil
	case "PushEvent":
		return &PushEvent{raw: e}, nil
	case "PullRequestEvent":
		return &PullRequestEvent{raw: e}, nil
	case "GollumEvent":
		return &GollumEvent{raw: e}, nil
	default:
		return nil, errors.New(fmt.Sprint("unknown event type: ", e.Type))
	}
}

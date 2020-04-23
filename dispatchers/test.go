package dispatchers

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type TestDispatcher struct {
	RepoName    string
	MessageSent string
	ShouldError bool
	MakeCalls   bool
	URL         string
}

func (td *TestDispatcher) Repo() string {
	return td.RepoName
}

func (td *TestDispatcher) SendMessage(message string, logger *logrus.Logger) error {
	if td.ShouldError {
		return errors.New("configured error")
	}
	td.MessageSent = message

	if td.MakeCalls {
		d := &SlackDispatcher{
			RepoName: "test",
			URL:      td.URL,
		}
		logger.WithField("event", "test_event_call").Info()
		return d.SendMessage(message, logger)
	}

	return nil
}

package dispatchers

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

type LocalDispatcher struct {
	RepoName string
	URL      string
}

func (ld *LocalDispatcher) Repo() string {
	return ld.RepoName
}

func (ld *LocalDispatcher) SendMessage(message string, logger *logrus.Logger) error {
	cmd := exec.Command("say", message)
	return cmd.Run()
}

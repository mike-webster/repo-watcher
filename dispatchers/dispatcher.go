package dispatchers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Dispatcher interface {
	Repo() string
	SendMessage(string, *logrus.Logger) error
}

type Dispatchers []Dispatcher

func (d *Dispatchers) ProcessMessage(repo string, message string, logger *logrus.Logger) error {
	for _, i := range *d {
		if strings.ToLower(i.Repo()) == strings.ToLower(repo) {
			return i.SendMessage(message, logger)
		}
	}

	return errors.New(fmt.Sprint("couldnt find dispatcher to match repo: ", repo))
}

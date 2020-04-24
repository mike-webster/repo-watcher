package dispatchers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type SlackDispatcher struct {
	RepoName string
	URL      string
}

func (sd *SlackDispatcher) Repo() string {
	return sd.RepoName
}

func (sd *SlackDispatcher) SendMessage(message string, logger *logrus.Logger) error {
	// these characters need to be escaped for slack
	// https://api.slack.com/reference/surfaces/formatting#escaping
	body := getBlockKitText(message)
	req, err := http.NewRequest("POST", sd.URL, strings.NewReader(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		logger.WithFields(logrus.Fields{
			"code": resp.StatusCode,
			"body": string(body),
			"url":  sd.URL,
		}).Error("non-200 response from extrnal call")

		return errors.New(fmt.Sprint("non-200 response: ", resp.StatusCode))
	}
	return nil
}

func getBlockKitText(text string) string {
	return fmt.Sprintf("{\"blocks\": [{\"type\": \"section\",\"text\": {\"type\": \"mrkdwn\", \"text\": \"%s\"}}]}", text)
}

package dispatchers

import (
	"encoding/json"
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
	body := getBlockKitText(message, logger)
	if len(body) < 1 {
		return errors.New("couldnt generate slack payload")
	}

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

func getBlockKitText(text string, logger *logrus.Logger) string {
	type slackText struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type slackBlock struct {
		Type string    `json:"type"`
		Text slackText `json:"text"`
	}
	type slackPayload struct {
		Blocks []slackBlock `json:"blocks"`
	}

	p := slackPayload{
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: slackText{
					Type: "mrkdwn",
					Text: text,
				},
			},
		},
	}

	bytes, err := json.Marshal(&p)
	if err != nil {
		logger.WithField("event", "cant_parse_slack_payload_to_json").Error(err)
		return ""
	}

	return string(bytes)
}

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type AppDependencies struct {
	logger      *logrus.Logger
	dispatchers Dispatchers
}

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

type SlackDispatcher struct {
	RepoName string
	URL      string
}

func (sd *SlackDispatcher) Repo() string {
	return sd.RepoName
}

func (sd *SlackDispatcher) SendMessage(message string, logger *logrus.Logger) error {
	// TODO: what are our options on formatting?
	body := fmt.Sprintf("{\"text\":\"%v\"}", message)
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

type TestDispatcher struct {
	RepoName    string
	MessageSent string
	ShouldError bool
}

func (td *TestDispatcher) Repo() string {
	return td.RepoName
}

func (td *TestDispatcher) SendMessage(message string, logger *logrus.Logger) error {
	if td.ShouldError {
		return errors.New("configured error")
	}

	td.MessageSent = message
	return nil
}

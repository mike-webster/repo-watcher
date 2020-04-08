package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	env "github.com/mike-webster/repo-watcher/env"
	models "github.com/mike-webster/repo-watcher/models"
)

// Log will print a message at the given level
func Log(message string, level string) {
	cfg := env.GetConfig()
	if cfg.LogLevel == "info" && level == "debug" {
		return
	}
	fmt.Println(time.Now().Format("2006-01-02T15:04:05-0700"), " -- ", strings.ToUpper(level), " -- ", message)
}

// MakeRequest will use the given url to make the appropriate request.
// If a string is provided for the id parameter, it will be included
// in the route.
// The return value should be able to be parsed into the desired struct
// as long as an error is not returned.
func MakeRequest(url string, id string, token string) (*[]byte, error) {
	Log("getting request", "debug")
	req, err := getRequest(url, id, token)
	if err != nil {
		return nil, err
	}

	Log("making request", "debug")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprint("non-200: ", resp.StatusCode))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()
	return &body, nil
}

func getRequest(url string, id string, token string) (*http.Request, error) {
	if len(id) > 0 {
		url = fmt.Sprintf(url, id)
	}

	Log("generating reqeust", "debug")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	Log("adding headers", "debug")
	req.Header.Add("Authorization", fmt.Sprint("token ", token))

	return req, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetPreviousIDs returns the last group of IDs returned from the events call
func GetPreviousIDs() (*[]string, error) {
	path := "history.txt"
	Log(fmt.Sprint("Looking in path: ", path, " for previous ids"), "debug")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			file, _ := os.Create("history.txt")
			file.Close()
			return &[]string{}, nil
		}

		return nil, err
	}

	ids := strings.Split(string(data), ",")
	return &ids, nil
}

// WriteNewIDs replaces the existing record with the current IDs
func WriteNewIDs(events []models.RepositoryEvent) error {
	path := "history.txt"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	ids := ""
	for _, event := range events {
		ids += fmt.Sprint(event.Raw().ID, ",")
	}
	if len(ids) > 0 {
		ids = strings.TrimSuffix(ids, ",")
	}

	_, err = w.WriteString(ids)
	if err != nil {
		return err
	}

	return w.Flush()
}

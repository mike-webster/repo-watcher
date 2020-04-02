package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Log will print a message at the given level
func Log(message string, level string) {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05-0700"), " -- ", strings.ToUpper(level), " -- ", message)
}

// MakeRequest will use the given key to make the appropriate request.
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

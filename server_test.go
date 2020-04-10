package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/mike-webster/repo-watcher/webhookmodels"
)

var server *Server

func TestMain(t *testing.T) {
	server = SetupServer("3199")

	testHealthcheck(t)
	testGithub(t)
}

func testHealthcheck(t *testing.T) {
	resp := performRequest(server.Engine, "GET", "/", map[string]string{}, []byte{})
	assert.Equal(t, CodeOK, resp.Code)
}

func testGithub(t *testing.T) {
	cases := []struct {
		Name            string
		Path            string
		Method          string
		ExpectedCode    int
		ExpectedMessage string
		Headers         map[string]string
		Body            interface{}
	}{
		{
			Name:            "TestEmptyHeaders",
			Path:            "/v1/github",
			Method:          "POST",
			ExpectedCode:    CodeInvalid,
			Headers:         map[string]string{},
			ExpectedMessage: `"{\"reason\":\"'ghRequestHeader.Event' Error:Field validation for 'Event' failed on the 'required' tag\",\"reason\":\"'ghRequestHeader.Secret' Error:Field validation for 'Secret' failed on the 'required' tag\"}"`,
		},
		{
			Name:            "TestEmptyHeaderSecret",
			Path:            "/v1/github",
			Method:          "POST",
			ExpectedCode:    CodeInvalid,
			Headers:         map[string]string{"X-GitHub-Event": "test"},
			ExpectedMessage: `"{\"reason\":\"'ghRequestHeader.Secret' Error:Field validation for 'Secret' failed on the 'required' tag\"}"`,
		},
		{
			Name:            "TestEmptyHeaderEvent",
			Path:            "/v1/github",
			Method:          "POST",
			ExpectedCode:    CodeInvalid,
			ExpectedMessage: `"{\"reason\":\"'ghRequestHeader.Event' Error:Field validation for 'Event' failed on the 'required' tag\"}"`,
			Headers:         map[string]string{"X-Hub-Signature": "test"},
		},
		{
			Name:            "TestEmptyBody",
			Path:            "/v1/github",
			Method:          "POST",
			ExpectedCode:    CodeInvalid,
			ExpectedMessage: `"{\"reason\":\"'PushEventPayload.Ref' Error:Field validation for 'Ref' failed on the 'required' tag\"}"`,
			Headers:         map[string]string{"X-GitHub-Event": "push", "X-Hub-Signature": "test"},
			Body:            webhookmodels.PushEventPayload{},
		},
		{
			Name:         "TestSuccess",
			Path:         "/v1/github",
			Method:       "POST",
			ExpectedCode: CodeNoContent,
			Headers:      map[string]string{"X-GitHub-Event": "push", "X-Hub-Signature": "push", "Content-Type": "application/json"},
			Body: webhookmodels.PushEventPayload{
				Ref: "test/ref",
				URL: "www.testurl.com",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			body, err := json.Marshal(c.Body)
			if err != nil {
				Log(fmt.Sprint("error parsing body: ", c.Body), "error")
				assert.Equal(t, nil, err)
			}
			resp := performRequest(server.Engine, c.Method, c.Path, c.Headers, body)
			assert.Equal(t, c.ExpectedCode, resp.Code, resp.Body.String())
			if len(c.ExpectedMessage) > 0 {
				assert.Equal(t, c.ExpectedMessage, resp.Body.String(), fmt.Sprintf("%v \n\t\t\t!=\n%v", c.ExpectedMessage, resp.Body.String()))
			}
		})
	}
}

func performRequest(r http.Handler, method string, path string, headers map[string]string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if len(body) > 0 {
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

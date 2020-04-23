package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gin-gonic/gin"
	dispatchers "github.com/mike-webster/repo-watcher/dispatchers"
	"github.com/mike-webster/repo-watcher/webhookmodels"
)

type testDeps struct {
	Router *gin.Engine
	Deps   *AppDependencies
}

func TestMain(t *testing.T) {
	deps := testSetup()

	testHealthcheck(t, deps)
	testGithub(t, deps)
	testParseEvent(t, deps)
}

func testSetup() *testDeps {
	deps := AppDependencies{
		logger: defaultLogger(),
		dispatchers: dispatchers.Dispatchers{
			&dispatchers.TestDispatcher{
				RepoName:  testWatch.Repo,
				URL:       testWatch.Webhook,
				MakeCalls: cfg.MakeTestCalls,
			},
		},
	}
	server := SetupServer("3199", &deps)
	return &testDeps{
		Router: server.Engine,
		Deps:   &deps,
	}
}

func testHealthcheck(t *testing.T, deps *testDeps) {
	resp := performRequest(deps.Router, "GET", "/", map[string]string{}, []byte{})
	assert.Equal(t, CodeOK, resp.Code)
}

func testGithub(t *testing.T, deps *testDeps) {
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
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
		},
	}

	t.Run("TestGitHub", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				body, err := json.Marshal(c.Body)
				if err != nil {
					assert.Equal(t, nil, err)
				}
				resp := performRequest(deps.Router, c.Method, c.Path, c.Headers, body)
				assert.Equal(t, c.ExpectedCode, resp.Code, resp.Body.String())
				if len(c.ExpectedMessage) > 0 {
					assert.Equal(t, c.ExpectedMessage, resp.Body.String(), fmt.Sprintf("%v \n\t\t\t!=\n%v", c.ExpectedMessage, resp.Body.String()))
				}
			})
		}
	})
}

func testParseEvent(t *testing.T, deps *testDeps) {
	cases := []struct {
		Name        string
		EventName   string
		Body        webhookmodels.Event
		Code        int
		DisplayName string
		Headers     map[string]string
	}{
		{
			Name:      "Create Event",
			EventName: "create",
			Body: &webhookmodels.PushEventPayload{
				Ref: "webby/test/ref",
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "create", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Gollum Event",
			EventName: "gollum",
			Body: &webhookmodels.GollumEventPayload{
				Pages: []webhookmodels.Page{
					webhookmodels.Page{
						Name: "test page name 1",
					},
					webhookmodels.Page{
						Name: "test page name 2",
					},
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "gollum", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Issue Comment Event",
			EventName: "issue_comment",
			Body: &webhookmodels.IssueCommentEventPayload{
				Action:  "commented",
				Comment: "test commment",
				Issue: webhookmodels.Issue{
					Title: "test title",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "issue_comment", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Issues Event",
			EventName: "issues",
			Body: &webhookmodels.IssuesEventPayload{
				Action: "opened",
				Issue: webhookmodels.Issue{
					Title: "test title",
					Body:  "test body",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "issues", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Project Card Event",
			EventName: "project_card",
			Body: &webhookmodels.ProjectCardEventPayload{
				Action: "created",
				Card: webhookmodels.Card{
					Note: "test note",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "project_card", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Project Column Event",
			EventName: "project_column",
			Body: &webhookmodels.ProjectColumnEventPayload{
				Action: "created",
				Column: webhookmodels.Column{
					Name: "test name",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "project_column", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Pull Request Event",
			EventName: "pull_request",
			Body: &webhookmodels.PullRequestEventPayload{
				Action: "opened",
				PullRequest: webhookmodels.PullRequest{
					Title: "test pull request title",
					Body:  "test pull request body",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "pull_request", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Pull Request Review Comment Event",
			EventName: "pull_request_review_comment",
			Body: &webhookmodels.PullRequestReviewCommentEventPayload{
				Action: "created",
				PullRequest: webhookmodels.PullRequest{
					Title: "test pull request title",
					State: "open",
				},
				Comment: webhookmodels.ReviewComment{
					Body: "test pull request body",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "pull_request_review_comment", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Pull Request Review Event",
			EventName: "pull_request_review",
			Body: &webhookmodels.PullRequestReviewEventPayload{
				Action: "created",
				PullRequest: webhookmodels.PullRequest{
					Title: "test pull request title",
					State: "open",
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "pull_request_review", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
		{
			Name:      "Push Event",
			EventName: "push",
			Body: &webhookmodels.PushEventPayload{
				Ref: "webby/test/ref",
				Commits: []interface{}{
					map[string]interface{}{
						"messge": "test commit 1",
					},
					map[string]interface{}{
						"message": "test commit 2",
					},
				},
				Sender: webhookmodels.User{
					Login: "mwebster",
				},
				Repo: webhookmodels.Repository{
					Name: "test",
				},
			},
			DisplayName: "Mike Webster",
			Code:        CodeNoContent,
			Headers:     map[string]string{"X-GitHub-Event": "push", "X-Hub-Signature": "push", "Content-Type": "application/json"},
		},
	}

	t.Run("VicariouslyTestParseEvent", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				var resp *httptest.ResponseRecorder
				if c.Body != nil {
					b, err := json.Marshal(c.Body)
					if err != nil {
						t.Error(err)
					}
					resp = performRequest(deps.Router, "POST", "/v1/github", c.Headers, b)
				} else {
					resp = performRequest(deps.Router, "POST", "/v1/github", c.Headers, nil)
				}
				assert.Equal(t, resp.Code, c.Code, resp.Body.String())
			})
		}
	})
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

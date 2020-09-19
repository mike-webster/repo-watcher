package webhookmodels

import (
	"fmt"
	"regexp"
	"strings"
)

func reworkLinks(body string) string {
	// looking for github links
	re := regexp.MustCompile(`\[([\w\s\d]+)\]\((https?:\/\/[\w\d./?=#_-]+)\)`)
	matched := re.FindAll([]byte(body), -1)
	if matched == nil {
		return body
	}
	for _, i := range matched {
		ghLink := string(i)
		reKey := regexp.MustCompile(`\[([\w\s\d]+)\]`)
		key := reKey.Find([]byte(ghLink))
		reVal := regexp.MustCompile(`\((https?:\/\/[\w\d./?=#_-]+)\)`)
		link := reVal.Find([]byte(ghLink))
		if len(key) > 0 && len(link) > 0 {
			body = strings.Replace(body, ghLink, fmt.Sprintf("<%s|%s>", link, key), -1)
		}
	}
	return body
}

func ShouldDeployMaster(e *Event) (*PullRequestEventPayload, bool) {
	local := *e
	prep, ok := local.(*PullRequestEventPayload)
	if !ok {
		// we only care if we're dealing with a pull request
		return nil, false
	}

	if prep.Action != "closed" {
		// we  don't want to deploy master on open, label, etc
		return nil, false
	}

	// is the repo subscribed?
	// ----> If not, return
	if "academy" != strings.ToLower(prep.Repo.Name) {
		return nil, false
	}

	return prep, true
}

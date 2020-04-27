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

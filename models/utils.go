package models

import (
	"strings"
)

func parseBranchName(ref string) string {
	sections := strings.Split(ref, "/")
	branch := sections[len(sections)-1]
	readable := strings.Split(branch, "-")
	return strings.Join(readable, " ")
}

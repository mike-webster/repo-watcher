package keys

type ContextKey string

var (
	TeamCityCreds ContextKey = "tc-creds"
	TeamCitySlackResponseURL ContextKey "slack-url"
)

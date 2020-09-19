package keys

import "context"

type ContextKey string

var (
	TeamCityCreds            ContextKey = "tc-creds"
	TeamCitySlackResponseURL ContextKey = "slack-url"
	AutoMerge                ContextKey = "auto-merge"
)

// Get will attempt  to retrieve the value from  the  context, first by using the
// key value. If that yields no results the  return value is the result of attempting
// to  retrieve the string value of the key.
// Note: This  is helpful for dealing  with gin contexts.
func Get(ctx context.Context, key ContextKey) interface{} {
	if ctx == nil {
		return nil
	}

	iVal := ctx.Value(key)
	if iVal != nil {
		return iVal
	}

	return ctx.Value(string(key))
}

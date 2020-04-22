package markdown

import "fmt"

// MarkdownLink takes a link and the text that should represent the link
// and returns the markdown representation
func MarkdownLink(url string, text string) string {
	return fmt.Sprintf("<%s|%s>", url, text)
}

// MarkdownBold takes a string and returns the bolded markdown equivalent
func MarkdownBold(text string) string {
	return fmt.Sprintf("*%s*", text)
}

// MarkdownItalic takes a string and returns the italicized markdown equivalent
func MarkdownItalic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

// MarkdownQuote takes a string and returns the quoted markdown equivalent
func MarkdownQuote(text string) string {
	// TODO: what will happen here if there are newlines?
	return fmt.Sprintf("> %s", text)
}

// MarkdownCode takes a string and returns the coded markdown equivalent
func MarkdownCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// MarkdownMultilineCode takes a string and returns the multi line coded markdown equivalent
func MarkdownMultilineCode(text string) string {
	return fmt.Sprintf("```%s```", text)
}

// MarkdownList takes a list of strings and returns a markdown equivalent
func MarkdownList(items []string) string {
	ret := ""
	for _, i := range items {
		ret += fmt.Sprintf("- %v\n", i)
	}
	return ret
}

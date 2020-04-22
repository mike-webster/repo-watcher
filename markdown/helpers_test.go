package markdown

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestMarkdownHelpers(t *testing.T) {
	text := "testing some text helpers"
	url := "www.google.com"
	t.Run("Link", func(t *testing.T) {
		expected := "<" + url + "|" + text + ">"
		assert.Equal(t, expected, MarkdownLink(url, text))
	})

	t.Run("Bold", func(t *testing.T) {
		expected := "*" + text + "*"
		assert.Equal(t, expected, MarkdownBold(text))
	})

	t.Run("Italic", func(t *testing.T) {
		expected := "_" + text + "_"
		assert.Equal(t, expected, MarkdownItalic(text))
	})

	t.Run("Quote", func(t *testing.T) {
		expected := "> " + text
		assert.Equal(t, expected, MarkdownQuote(text))
	})

	t.Run("Code", func(t *testing.T) {
		expected := "`" + text + "`"
		assert.Equal(t, expected, MarkdownCode(text))
	})

	t.Run("MultilineCode", func(t *testing.T) {
		expected := "```" + text + "```"
		assert.Equal(t, expected, MarkdownMultilineCode(text))
	})

	t.Run("List", func(t *testing.T) {
		expected := "- Testing1\n- Testing2\n- Testing3\n"
		assert.Equal(t, expected, MarkdownList([]string{"Testing1", "Testing2", "Testing3"}))
	})
}

package main

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/mike-webster/repo-watcher/markdown"
)

func TestMarkdownInSlack(t *testing.T) {
	deps := testSetup()

	t.Run("Link", func(t *testing.T) {
		text := "testing markdown"
		url := "http://www.github.com/mike-webster/repo-watcher"
		message := markdown.MarkdownLink(url, text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("Bold", func(t *testing.T) {
		text := "testing markdown"
		message := markdown.MarkdownBold(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("Italic", func(t *testing.T) {
		text := "testing markdown"
		message := markdown.MarkdownItalic(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("Quote", func(t *testing.T) {
		text := "testing markdown"
		message := markdown.MarkdownQuote(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("Code", func(t *testing.T) {
		text := "testing markdown"
		message := markdown.MarkdownCode(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("MultiCode", func(t *testing.T) {
		text := "testing markdown"
		message := markdown.MarkdownMultilineCode(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", message, deps.Deps.logger))
	})
	t.Run("List", func(t *testing.T) {
		text := []string{"test1", "test2", "test3"}
		list := markdown.MarkdownList(text)
		assert.Equal(t, nil, deps.Deps.dispatchers.ProcessMessage("test", list, deps.Deps.logger))
	})
}

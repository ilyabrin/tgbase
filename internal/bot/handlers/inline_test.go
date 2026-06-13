package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

type inlineCtx struct {
	telebot.Context
	query    *telebot.Query
	answered *telebot.QueryResponse
}

func (c *inlineCtx) Query() *telebot.Query { return c.query }
func (c *inlineCtx) Answer(r *telebot.QueryResponse) error {
	c.answered = r
	return nil
}

func TestInlineHandler_WithText(t *testing.T) {
	h := InlineHandler()
	c := &inlineCtx{query: &telebot.Query{Text: "hello"}}

	require.NoError(t, h(c))
	require.NotNil(t, c.answered)
	assert.Len(t, c.answered.Results, 1)

	result := c.answered.Results[0].(*telebot.ArticleResult)
	assert.Equal(t, "hello", result.Title)
	assert.Equal(t, "hello", result.Text)
	assert.Equal(t, 60, c.answered.CacheTime)
}

func TestInlineHandler_EmptyQuery_UsesFallback(t *testing.T) {
	h := InlineHandler()
	c := &inlineCtx{query: &telebot.Query{Text: "   "}}

	require.NoError(t, h(c))
	require.NotNil(t, c.answered)
	result := c.answered.Results[0].(*telebot.ArticleResult)
	assert.Equal(t, "...", result.Title)
}

func TestInlineHandler_NilQuery(t *testing.T) {
	h := InlineHandler()
	c := &inlineCtx{query: nil}
	require.NoError(t, h(c))
	assert.Nil(t, c.answered)
}

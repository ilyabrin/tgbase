package handlers

import (
	"strings"

	"gopkg.in/telebot.v3"
)

// InlineHandler handles inline queries (when users type "@botname <query>" in
// any chat). Returns article results that the user can tap to send.
//
// Register with:
//
//	b.Handle(telebot.OnQuery, handlers.InlineHandler())
//
// Customise the results slice to fit your use case (search your DB, etc.).
func InlineHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		query := c.Query()
		if query == nil {
			return nil
		}

		text := strings.TrimSpace(query.Text)
		if text == "" {
			text = "..."
		}

		results := telebot.Results{
			&telebot.ArticleResult{
				ResultBase:  telebot.ResultBase{ID: "result_0"},
				Title:       text,
				Description: "Send «" + text + "» to the chat",
				Text:        text,
			},
		}

		return c.Answer(&telebot.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})
	}
}

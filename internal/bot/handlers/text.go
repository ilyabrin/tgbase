package handlers

import (
	"tgbase/internal/i18n"

	"gopkg.in/telebot.v3"
)

// TextHandler handles text messages
func TextHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en" // Fallback to English
		}
		message := i18n.Localize(lang, "echo", map[string]any{
			"Text": c.Text(),
		})
		if message == "" {
			message = "You said: " + c.Text()
		}
		return c.Send(message)
	}
}
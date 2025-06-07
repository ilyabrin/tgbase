package handlers

import (
	"tgbase/internal/i18n"

	"gopkg.in/telebot.v3"
)

// StartHandler handles the /start command
func StartHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en" // Fallback to English
		}
		message := i18n.Localize(lang, "welcome", nil)
		if message == "" {
			message = "Welcome to the bot!"
		}
		return c.Send(message)
	}
}
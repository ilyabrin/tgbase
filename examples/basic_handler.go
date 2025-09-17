// Package examples provides examples of how to implement common bot handlers
package examples

import (
	"tgbase/internal/i18n"

	"gopkg.in/telebot.v3"
)

// BasicHandler demonstrates a simple command handler that responds with a static message
func BasicHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en" // Fallback to English
		}

		// Try to get localized message, fallback to English text
		message := i18n.Localize(lang, "basic_response", nil)
		if message == "" {
			message = "This is a basic handler response!"
		}

		return c.Send(message)
	}
}

// EchoHandler demonstrates how to echo user input with additional processing
func EchoHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userText := c.Text()
		if userText == "" {
			return c.Send("Please send me some text to echo!")
		}

		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en"
		}

		// Use localization with template data
		message := i18n.Localize(lang, "echo_response", map[string]any{
			"Text":     userText,
			"Username": c.Sender().FirstName,
		})

		if message == "" {
			message = "Hello " + c.Sender().FirstName + "! You said: " + userText
		}

		return c.Send(message)
	}
}
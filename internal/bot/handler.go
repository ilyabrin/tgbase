package bot

import (
	"gopkg.in/telebot.v3"
)

func (b *Bot) registerHandlers() {
	// Обработка команды /start
	b.bot.Handle("/start", func(c telebot.Context) error {
		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en" // Fallback на английский
		}
		message := b.i18n.Localize(lang, "welcome", nil)
		return c.Send(message)
	})

	// Обработка текстовых сообщений
	b.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en" // Fallback на английский
		}
		message := b.i18n.Localize(lang, "echo", map[string]any{
			"Text": c.Text(),
		})
		return c.Send(message)
	})
}

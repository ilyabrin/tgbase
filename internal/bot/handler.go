package bot

import (
	"gopkg.in/telebot.v3"
)

func (b *Bot) registerHandlers() {
	// Обработка команды /start
	b.bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Welcome to the bot!")
	})

	// Обработка текстовых сообщений
	b.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return c.Send("Echo: " + c.Text())
	})
}

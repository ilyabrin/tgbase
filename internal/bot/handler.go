package bot

import (
	"tgbase/internal/bot/handlers"

	"gopkg.in/telebot.v3"
)

func (b *Bot) registerHandlers() {
	// Register /start command handler
	b.bot.Handle("/start", handlers.StartHandler(b.i18n))

	// Register text message handler
	b.bot.Handle(telebot.OnText, handlers.TextHandler(b.i18n))
}
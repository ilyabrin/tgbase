package bot

import (
	"tgbase/internal/bot/handlers"

	"gopkg.in/telebot.v3"
)

// RegisterHandlers wires all application handlers to the bot.
// Add your handlers here.
func (b *Bot) RegisterHandlers() {
	b.Handle("/start", handlers.StartHandler(b.i18n))
	b.Handle(telebot.OnText, handlers.TextHandler(b.i18n))

	if b.redis != nil {
		b.Handle("/redis2", handlers.Redis2Handler(b.redis))
		b.Handle(handlers.BtnToggle, handlers.HandleRedis2Button(b.redis))
	}
}

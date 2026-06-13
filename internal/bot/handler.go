package bot

import (
	"time"

	"tgbase/internal/bot/handlers"
	"tgbase/internal/fsm"
	"tgbase/pkg/middleware"

	"gopkg.in/telebot.v3"
)

// RegisterHandlers wires all application handlers to the bot.
// Add your handlers here.
func (b *Bot) RegisterHandlers() {
	b.Use(middleware.Logger(b.logger))
	b.Use(middleware.RateLimit(30, time.Minute))

	b.Handle("/start", handlers.StartHandler(b.i18n))
	b.Handle(telebot.OnQuery, handlers.InlineHandler())

	if b.redis != nil {
		b.Handle("/redis2", handlers.Redis2Handler(b.redis))
		b.Handle(handlers.BtnToggle, handlers.HandleRedis2Button(b.redis))

		// FSM: multi-step registration demo.
		// TextHandler is the fallback for users with no active state.
		f := fsm.New(fsm.NewRedisStorage(b.redis, fsm.WithTTL(3600))).
			Fallback(handlers.TextHandler(b.i18n))

		b.Handle("/register", handlers.RegisterStart(f))
		b.Handle(telebot.OnText, f.Route(
			fsm.On(handlers.StateAskName, handlers.RegisterAskName(f)),
			fsm.On(handlers.StateAskAge, handlers.RegisterAskAge(f)),
		))
	} else {
		b.Handle(telebot.OnText, handlers.TextHandler(b.i18n))
	}

	// Telegram Stars payments.
	// Define your products and uncomment:
	//
	// product := handlers.StarProduct{
	//     Title:       "Premium",
	//     Description: "Unlock all features for 30 days",
	//     Payload:     "premium_1month",
	//     Stars:       100,
	// }
	// b.Handle("/buy",              handlers.SendInvoice(product))
	// b.Handle(telebot.OnCheckout,  handlers.PreCheckout())
	// b.Handle(telebot.OnPayment,   handlers.PaymentSuccess(b.db))

	// Admin-only example - uncomment and add your handler:
	// b.Handle("/admin", adminHandler, middleware.AdminOnly(b.adminIDs))
}

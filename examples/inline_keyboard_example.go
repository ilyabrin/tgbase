// Package examples shows how to create interactive inline keyboards
package examples

import (
	"context"
	"fmt"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"

	"gopkg.in/telebot.v3"
)

const (
	// Button callback data constants
	BtnYes    = "btn_yes"
	BtnNo     = "btn_no"
	BtnHelp   = "btn_help"
	BtnMenu   = "btn_menu"
	BtnOption = "btn_option_"
)

// MenuHandler demonstrates how to create an interactive menu with inline buttons
func MenuHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Create inline keyboard
		inlineKeys := &telebot.ReplyMarkup{}

		btnHelp := inlineKeys.Data("📖 Help", BtnHelp)
		btnStats := inlineKeys.Data("📊 Statistics", "btn_stats")
		btnSettings := inlineKeys.Data("⚙️ Settings", "btn_settings")

		// Arrange buttons in rows
		inlineKeys.Inline(
			inlineKeys.Row(btnHelp, btnStats),
			inlineKeys.Row(btnSettings),
		)

		lang := c.Sender().LanguageCode
		if lang == "" {
			lang = "en"
		}

		message := i18n.Localize(lang, "menu_message", nil)
		if message == "" {
			message = "🤖 Bot Menu\n\nChoose an option:"
		}

		return c.Send(message, inlineKeys)
	}
}

// QuestionHandler demonstrates yes/no questions with callbacks
func QuestionHandler(i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		inlineKeys := &telebot.ReplyMarkup{}

		btnYes := inlineKeys.Data("✅ Yes", BtnYes)
		btnNo := inlineKeys.Data("❌ No", BtnNo)

		inlineKeys.Inline(
			inlineKeys.Row(btnYes, btnNo),
		)

		message := "Do you like this bot?"
		return c.Send(message, inlineKeys)
	}
}

// HandleCallbacks demonstrates how to handle button callbacks
func HandleCallbacks(redis *redis.Client, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Answer the callback to remove loading state
		defer c.Respond()

		callback := c.Callback()
		userID := c.Sender().ID

		var responseMessage string

		switch callback.Data {
		case BtnYes:
			responseMessage = "Great! Thank you for the positive feedback! 😊"
			// Store user preference in Redis
			if redis != nil {
				ctx := context.Background()
				key := fmt.Sprintf("user:%d:likes_bot", userID)
				redis.Set(ctx, key, "yes", 3600) // Store for 1 hour
			}

		case BtnNo:
			responseMessage = "Thanks for the feedback. How can we improve? 🤔"
			if redis != nil {
				ctx := context.Background()
				key := fmt.Sprintf("user:%d:likes_bot", userID)
				redis.Set(ctx, key, "no", 3600)
			}

		case BtnHelp:
			responseMessage = "📖 Bot Help\n\n" +
				"Available commands:\n" +
				"/start - Start the bot\n" +
				"/menu - Show main menu\n" +
				"/question - Ask a question\n" +
				"/redis2 - Redis demo (if available)"

		case "btn_stats":
			responseMessage = "📊 Bot Statistics\n\n" +
				"Total users: 42\n" +
				"Messages processed: 1337\n" +
				"Uptime: 24h 30m"

		case "btn_settings":
			responseMessage = "⚙️ Settings\n\n" +
				"Language: Auto-detect\n" +
				"Notifications: Enabled\n" +
				"Theme: Default"

		default:
			responseMessage = "Unknown button pressed."
		}

		// Edit the original message with the response
		return c.Edit(responseMessage)
	}
}

// DynamicMenuHandler demonstrates a menu that changes based on user state
func DynamicMenuHandler(redis *redis.Client, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		userID := c.Sender().ID

		inlineKeys := &telebot.ReplyMarkup{}
		var buttons []telebot.Row

		// Check if user has premium access (example)
		var isPremium bool
		if redis != nil {
			premiumKey := fmt.Sprintf("user:%d:premium", userID)
			premiumStatus, _ := redis.Get(ctx, premiumKey)
			isPremium = premiumStatus == "true"
		}

		// Basic buttons for all users
		btnBasic := inlineKeys.Data("🏠 Basic Features", "btn_basic")
		buttons = append(buttons, inlineKeys.Row(btnBasic))

		// Premium-only buttons
		if isPremium {
			btnPremium := inlineKeys.Data("⭐ Premium Features", "btn_premium")
			buttons = append(buttons, inlineKeys.Row(btnPremium))
		} else {
			btnUpgrade := inlineKeys.Data("💎 Upgrade to Premium", "btn_upgrade")
			buttons = append(buttons, inlineKeys.Row(btnUpgrade))
		}

		// Settings button for everyone
		btnSettings := inlineKeys.Data("⚙️ Settings", "btn_settings")
		buttons = append(buttons, inlineKeys.Row(btnSettings))

		inlineKeys.Inline(buttons...)

		message := "🎯 Dynamic Menu\n\n"
		if isPremium {
			message += "Welcome, Premium User! 👑"
		} else {
			message += "Welcome! Upgrade for more features."
		}

		return c.Send(message, inlineKeys)
	}
}
package handlers

import (
	"context"
	"fmt"
	"tgbase/internal/redis"

	"gopkg.in/telebot.v3"
)

const (
	// Redis2Key is the key prefix for the redis2 command state in Redis
	Redis2Key = "redis2"
	// BtnToggle is the data for the toggle button
	BtnToggle = "btn_toggle"
)

// Redis2Handler handles the /redis2 command
func Redis2Handler(redis *redis.Client) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Create unique key for this user in this chat
		key := fmt.Sprintf("chat:%d:user:%d:%s", c.Chat().ID, c.Sender().ID, Redis2Key)

		// Initialize state in Redis
		if err := redis.Del(context.Background(), key); err != nil {
			return fmt.Errorf("failed to initialize state: %w", err)
		}

		// Create toggle button
		inlineKeys := &telebot.ReplyMarkup{}
		btnToggle := inlineKeys.Data("Click me!", BtnToggle)
		inlineKeys.Inline(
			inlineKeys.Row(btnToggle),
		)

		// Send initial message
		return c.Send("Toggle state: unpressed", inlineKeys)
	}
}

// HandleRedis2Button handles callback queries from redis2 toggle button
func HandleRedis2Button(redis *redis.Client) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()

		// Create unique key for this user in this chat
		key := fmt.Sprintf("chat:%d:user:%d:%s", c.Chat().ID, c.Sender().ID, Redis2Key)

		// Get current state
		exists, err := redis.Exists(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to check state: %w", err)
		}

		var messageText string
		if !exists {
			// Set pressed state
			if err := redis.Set(ctx, key, "pressed", 0); err != nil {
				return fmt.Errorf("failed to set state: %w", err)
			}
			messageText = fmt.Sprintf("Toggle state: pressed: %s", c.Sender().FirstName)
		} else {
			// Remove pressed state
			if err := redis.Del(ctx, key); err != nil {
				return fmt.Errorf("failed to remove state: %w", err)
			}
			messageText = "Toggle state: unpressed"
		}

		// Answer callback query to remove loading state
		if err := c.Respond(); err != nil {
			return fmt.Errorf("failed to respond to callback: %w", err)
		}

		// Recreate the inline keyboard
		inlineKeys := &telebot.ReplyMarkup{}
		btnToggle := inlineKeys.Data("Click me!", BtnToggle)
		inlineKeys.Inline(
			inlineKeys.Row(btnToggle),
		)

		// Update message text with the inline keyboard
		return c.Edit(messageText, inlineKeys)
	}
}

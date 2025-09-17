// Package examples shows how to use Redis for caching and session management
package examples

import (
	"context"
	"fmt"
	"strconv"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"
	"time"

	"gopkg.in/telebot.v3"
)

// CacheHandler demonstrates how to use Redis for caching
func CacheHandler(redis *redis.Client, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		userID := c.Sender().ID

		// Create a cache key for this user
		cacheKey := fmt.Sprintf("user_cache:%d", userID)

		// Try to get cached data
		cachedValue, err := redis.Get(ctx, cacheKey)
		if err != nil {
			return fmt.Errorf("failed to get cached data: %w", err)
		}

		var message string
		if cachedValue != "" {
			message = fmt.Sprintf("Cached data: %s", cachedValue)
		} else {
			// Simulate expensive operation
			expensiveData := fmt.Sprintf("Computed at %s", time.Now().Format("15:04:05"))

			// Cache for 30 seconds
			if err := redis.Set(ctx, cacheKey, expensiveData, 30); err != nil {
				return fmt.Errorf("failed to cache data: %w", err)
			}

			message = fmt.Sprintf("New data cached: %s", expensiveData)
		}

		return c.Send(message)
	}
}

// SessionHandler demonstrates session management with Redis
func SessionHandler(redis *redis.Client, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		userID := c.Sender().ID

		// Session key for this user
		sessionKey := fmt.Sprintf("session:%d", userID)

		// Get current session data
		sessionData, err := redis.HGetAll(ctx, sessionKey)
		if err != nil {
			return fmt.Errorf("failed to get session data: %w", err)
		}

		// Update session with new interaction
		if err := redis.HSet(ctx, sessionKey, "last_command", "/session"); err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}

		if err := redis.HSet(ctx, sessionKey, "last_interaction", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}

		// Increment interaction counter
		countStr := sessionData["interaction_count"]
		count := 0
		if countStr != "" {
			count, _ = strconv.Atoi(countStr)
		}
		count++

		if err := redis.HSet(ctx, sessionKey, "interaction_count", strconv.Itoa(count)); err != nil {
			return fmt.Errorf("failed to update interaction count: %w", err)
		}

		message := fmt.Sprintf("Session Info:\n"+
			"Interactions: %d\n"+
			"Last Command: %s\n"+
			"Last Update: %s",
			count,
			sessionData["last_command"],
			sessionData["last_interaction"])

		return c.Send(message)
	}
}

// CounterHandler demonstrates a simple shared counter using Redis
func CounterHandler(redis *redis.Client, i18n *i18n.I18n) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		counterKey := "global_counter"

		// Increment the global counter
		newValue, err := redis.Incr(ctx, counterKey)
		if err != nil {
			return fmt.Errorf("failed to increment counter: %w", err)
		}

		message := fmt.Sprintf("Global counter: %d", newValue)
		return c.Send(message)
	}
}
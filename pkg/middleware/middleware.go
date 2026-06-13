// Package middleware provides reusable telebot middleware.
//
// Register globally:
//
//	b.Use(middleware.Logger(logger))
//
// Or per-handler:
//
//	b.Handle("/ban", banHandler, middleware.AdminOnly(cfg.AdminIDs))
package middleware

import (
	"fmt"
	"sync"
	"time"

	"tgbase/pkg/logger"

	"gopkg.in/telebot.v3"
)

// AdminOnly rejects requests from users not in the allowlist.
// Non-admins receive a silent drop (no response); pass onReject to customise.
func AdminOnly(adminIDs []int64, onReject ...telebot.HandlerFunc) telebot.MiddlewareFunc {
	allowed := make(map[int64]struct{}, len(adminIDs))
	for _, id := range adminIDs {
		allowed[id] = struct{}{}
	}

	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if _, ok := allowed[c.Sender().ID]; ok {
				return next(c)
			}
			if len(onReject) > 0 {
				return onReject[0](c)
			}
			return nil
		}
	}
}

// Logger logs every incoming update: user, chat, and message/command text.
func Logger(l *logger.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			sender := c.Sender()
			name := sender.FirstName
			if sender.LastName != "" {
				name += " " + sender.LastName
			}

			var what string
			switch {
			case c.Message() != nil && c.Message().Text != "":
				what = fmt.Sprintf("%q", c.Message().Text)
			case c.Callback() != nil:
				what = fmt.Sprintf("callback %q", c.Callback().Data)
			default:
				what = "update"
			}

			l.Info("user %d (%s): %s", sender.ID, name, what)
			return next(c)
		}
	}
}

// RateLimit allows at most n messages per user per window duration.
// Excess messages are silently dropped unless onReject is provided.
func RateLimit(n int, per time.Duration, onReject ...telebot.HandlerFunc) telebot.MiddlewareFunc {
	type entry struct {
		mu    sync.Mutex
		count int
		reset time.Time
	}

	var buckets sync.Map

	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			id := c.Sender().ID

			val, _ := buckets.LoadOrStore(id, &entry{reset: time.Now().Add(per)})
			e := val.(*entry)

			e.mu.Lock()
			now := time.Now()
			if now.After(e.reset) {
				e.count = 0
				e.reset = now.Add(per)
			}
			e.count++
			allowed := e.count <= n
			e.mu.Unlock()

			if allowed {
				return next(c)
			}
			if len(onReject) > 0 {
				return onReject[0](c)
			}
			return nil
		}
	}
}

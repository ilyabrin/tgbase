package bot

import (
	"context"
	"time"

	"tgbase/pkg/logger"

	"gopkg.in/telebot.v3"
)

// BroadcastResult summarises a completed broadcast.
type BroadcastResult struct {
	Sent   int
	Failed int
	Errors []BroadcastError
}

// BroadcastError pairs a user ID with the send error.
type BroadcastError struct {
	UserID int64
	Err    error
}

// BroadcastOption configures a Broadcast call.
type BroadcastOption func(*broadcastCfg)

type broadcastCfg struct {
	delay   time.Duration
	onError func(int64, error)
}

// WithBroadcastDelay sets the pause between individual sends.
// Default is 50 ms (≈ 20 msg/s), safely below Telegram's 30 msg/s limit.
func WithBroadcastDelay(d time.Duration) BroadcastOption {
	return func(c *broadcastCfg) { c.delay = d }
}

// WithBroadcastOnError registers a callback invoked for each failed send.
func WithBroadcastOnError(fn func(userID int64, err error)) BroadcastOption {
	return func(c *broadcastCfg) { c.onError = fn }
}

// sender is the subset of *telebot.Bot used by broadcast - allows mocking in tests.
type sender interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

// Broadcast sends what to every user in userIDs, respecting ctx cancellation.
// Failed sends are collected in the result; the broadcast always completes
// (or stops early on context cancellation) without returning an error itself.
func (b *Bot) Broadcast(ctx context.Context, userIDs []int64, what interface{}, opts ...BroadcastOption) BroadcastResult {
	cfg := &broadcastCfg{delay: 50 * time.Millisecond}
	for _, o := range opts {
		o(cfg)
	}
	return runBroadcast(ctx, b.bot, b.logger, userIDs, what, cfg)
}

func runBroadcast(ctx context.Context, s sender, l *logger.Logger, userIDs []int64, what interface{}, cfg *broadcastCfg) BroadcastResult {
	var result BroadcastResult

	for _, id := range userIDs {
		if ctx.Err() != nil {
			return result
		}

		_, err := s.Send(telebot.ChatID(id), what)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, BroadcastError{UserID: id, Err: err})
			if cfg.onError != nil {
				cfg.onError(id, err)
			}
			l.Error("broadcast: user %d: %v", id, err)
		} else {
			result.Sent++
		}

		if cfg.delay > 0 {
			select {
			case <-ctx.Done():
				return result
			case <-time.After(cfg.delay):
			}
		}
	}

	return result
}

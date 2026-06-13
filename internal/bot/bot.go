package bot

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"tgbase/internal/database"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"
	"tgbase/pkg/logger"

	"gopkg.in/telebot.v3"
)

// Bot wraps telebot.Bot with application dependencies.
type Bot struct {
	bot           *telebot.Bot
	db            database.Database
	redis         *redis.Client
	i18n          *i18n.I18n
	logger        *logger.Logger
	webhookAddr   string
	pollerTimeout time.Duration
}

// Option configures a Bot.
type Option func(*Bot)

func WithDB(db database.Database) Option {
	return func(b *Bot) { b.db = db }
}

func WithRedis(r *redis.Client) Option {
	return func(b *Bot) { b.redis = r }
}

func WithI18n(i *i18n.I18n) Option {
	return func(b *Bot) { b.i18n = i }
}

func WithLogger(l *logger.Logger) Option {
	return func(b *Bot) { b.logger = l }
}

// WithWebhook switches the bot from long-polling to webhook mode.
func WithWebhook(listenAddr string) Option {
	return func(b *Bot) { b.webhookAddr = listenAddr }
}

// WithPollerTimeout overrides the default 10s long-poll timeout.
func WithPollerTimeout(d time.Duration) Option {
	return func(b *Bot) { b.pollerTimeout = d }
}

// New creates a bot. All options are optional; sensible defaults are applied.
// Register handlers after construction, then call Run or Start.
func New(token string, opts ...Option) (*Bot, error) {
	b := &Bot{pollerTimeout: 10 * time.Second}
	for _, o := range opts {
		o(b)
	}
	if b.logger == nil {
		b.logger = logger.NewLogger()
	}

	var poller telebot.Poller
	if b.webhookAddr != "" {
		poller = &telebot.Webhook{Listen: b.webhookAddr}
	} else {
		poller = &telebot.LongPoller{Timeout: b.pollerTimeout}
	}

	tgb, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: poller,
	})
	if err != nil {
		return nil, err
	}
	b.bot = tgb
	return b, nil
}

// Handle registers an endpoint handler, optionally with middleware.
// Returns the bot for fluent chaining.
func (b *Bot) Handle(endpoint any, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) *Bot {
	b.bot.Handle(endpoint, h, m...)
	return b
}

// Use registers global middleware applied to all handlers.
func (b *Bot) Use(middleware ...telebot.MiddlewareFunc) *Bot {
	b.bot.Use(middleware...)
	return b
}

// Run starts the bot and blocks until SIGINT or SIGTERM is received.
func (b *Bot) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	_ = b.Start(ctx)
}

// Start runs the bot until ctx is cancelled.
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("starting telegram bot...")
	go b.bot.Start()
	<-ctx.Done()
	b.bot.Stop()
	b.logger.Info("telegram bot stopped")
	return nil
}

package bot

import (
	"context"
	"time"

	"tgbase/internal/database"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"
	"tgbase/pkg/logger"

	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot    *telebot.Bot
	db     database.Database
	redis  *redis.Client
	i18n   *i18n.I18n
	logger *logger.Logger
}

func NewBot(token string, db database.Database, redis *redis.Client, i18n *i18n.I18n, logger *logger.Logger) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	b := &Bot{
		bot:    bot,
		db:     db,
		redis:  redis,
		i18n:   i18n,
		logger: logger,
	}

	// Регистрация обработчиков
	b.registerHandlers()

	return b, nil
}

// NewBotWebhook создает нового бота с использованием вебхуков
func NewBotWebhook(token string, db database.Database, redis *redis.Client, i18n *i18n.I18n, logger *logger.Logger) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.Webhook{Listen: ":8080"},
	})
	if err != nil {
		return nil, err
	}

	b := &Bot{
		bot:    bot,
		db:     db,
		redis:  redis,
		i18n:   i18n,
		logger: logger,
	}

	// Регистрация обработчиков
	b.registerHandlers()

	return b, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("starting telegram bot...")

	// Запуск бота в отдельной горутине
	go b.bot.Start()

	// Ожидание сигнала завершения
	<-ctx.Done()
	b.bot.Stop()
	b.logger.Info("telegram bot stopped")
	return nil
}

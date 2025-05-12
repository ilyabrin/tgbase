package bot

import (
	"context"
	"time"

	"tgbase/internal/database"
	"tgbase/internal/redis"
	"tgbase/pkg/logger"

	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot    *telebot.Bot
	db     database.Database
	redis  *redis.Client
	logger *logger.Logger
}

func NewBot(token string, db database.Database, redis *redis.Client, logger *logger.Logger) (*Bot, error) {
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

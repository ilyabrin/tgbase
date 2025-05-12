package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"tgbase/config"
	"tgbase/internal/bot"
	"tgbase/internal/database"
	"tgbase/internal/redis"
	"tgbase/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера
	logger := logger.NewLogger()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatal("failed to load config: ", err)
	}

	// Инициализация базы данных (Postgres или SQLite)
	var db database.Database
	if cfg.Database.Type == "postgres" {
		db, err = database.NewPostgresDB(ctx, cfg.Database.Postgres.DSN)
	} else {
		db, err = database.NewSQLiteDB(ctx, cfg.Database.SQLite.Path)
	}
	if err != nil {
		logger.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	// Инициализация Redis (опционально)
	var redisClient *redis.Client
	if cfg.Redis.Enabled {
		redisClient, err = redis.NewRedisClient(ctx, struct {
			Addr     string
			Password string
			DB       int
		}{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err != nil {
			logger.Fatal("failed to connect to redis: ", err)
		}
		defer redisClient.Close()
	}

	// Инициализация Telegram бота
	tgBot, err := bot.NewBot(cfg.Telegram.Token, db, redisClient, logger)
	if err != nil {
		logger.Fatal("failed to initialize bot: ", err)
	}

	// Запуск бота
	go func() {
		if err := tgBot.Start(ctx); err != nil {
			logger.Fatal("bot stopped: ", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down application...")
}

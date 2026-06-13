package main

import (
	"context"
	"log"

	"tgbase/config"
	"tgbase/internal/bot"
	"tgbase/internal/database"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	ctx := context.Background()

	db, err := database.FromConfig(ctx, cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	locale, err := i18n.NewI18n("locales")
	if err != nil {
		log.Fatal("failed to initialize i18n: ", err)
	}

	opts := []bot.Option{
		bot.WithDB(db),
		bot.WithI18n(locale),
	}

	if cfg.Redis.Enabled {
		redisClient, err := redis.NewRedisClient(ctx, redis.Config{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err != nil {
			log.Fatal("failed to connect to redis: ", err)
		}
		defer redisClient.Close()
		opts = append(opts, bot.WithRedis(redisClient))
	}

	b, err := bot.New(cfg.Telegram.Token, opts...)
	if err != nil {
		log.Fatal("failed to initialize bot: ", err)
	}

	b.RegisterHandlers()
	b.Run()
}

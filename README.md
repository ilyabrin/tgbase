# tgbase - Go Telegram Bot Template

[![CI Status](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml)
[![Redis Tests](https://github.com/ilyabrin/tgbase/actions/workflows/redis-tests.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/redis-tests.yml)
[![Codecov](https://codecov.io/gh/ilyabrin/tgbase/branch/main/graph/badge.svg)](https://codecov.io/gh/ilyabrin/tgbase)
[![Go Report Card](https://goreportcard.com/badge/github.com/ilyabrin/tgbase)](https://goreportcard.com/report/github.com/ilyabrin/tgbase)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Developer-friendly template for building Telegram bots in Go. Clone, configure, add handlers.

**Batteries included:** database (Postgres/SQLite), Redis, FSM, middleware, i18n, payments (Stars), Docker.

---

## Quick start

```bash
git clone https://github.com/ilyabrin/tgbase && cd tgbase
go mod tidy
# set your token in config.yaml
go run cmd/app/main.go
```

---

## Project structure

```text
tgbase/
├── cmd/app/main.go               # entry point - wires everything together
├── config/
│   └── config.go                 # YAML config + env overrides
├── internal/
│   ├── bot/
│   │   ├── bot.go                # Bot type, functional options, Run/Start
│   │   ├── handler.go            # RegisterHandlers - add your handlers here
│   │   └── handlers/
│   │       ├── start.go          # /start
│   │       ├── text.go           # OnText fallback
│   │       ├── redis2.go         # Redis toggle demo
│   │       ├── registration.go   # FSM multi-step flow demo
│   │       └── payment.go        # Telegram Stars payments
│   ├── database/
│   │   ├── database.go           # Database + SoftDeleteDatabase interfaces
│   │   ├── postgres.go           # PostgreSQL implementation
│   │   └── sqlite.go             # SQLite implementation
│   ├── fsm/
│   │   ├── fsm.go                # FSM: New, Route, On, SetState, GetData…
│   │   └── storage.go            # Storage interface, RedisStorage, MemoryStorage
│   ├── i18n/
│   │   └── i18n.go               # go-i18n wrapper, locales/*.yaml
│   └── redis/
│       ├── redis.go              # Client, Config, NewRedisClient, NewMockClient
│       └── real_redis.go         # go-redis/v9 implementation
├── pkg/
│   ├── logger/logger.go          # stdlib logger wrapper
│   └── middleware/middleware.go  # AdminOnly, Logger, RateLimit
├── config.yaml                   # application configuration
├── docker-compose.yml
└── Dockerfile
```

---

## Configuration (`config.yaml`)

```yaml
database:
  type: sqlite          # or "postgres"
  postgres:
    dsn: "host=localhost user=postgres password=secret dbname=app port=5432 sslmode=disable"
  sqlite:
    path: "app.db"

redis:
  enabled: true
  addr: "localhost:6379"
  password: ""
  db: 0

telegram:
  token: "YOUR_BOT_TOKEN"
  admin_ids:
    - 123456789          # your Telegram user ID
```

Environment variable overrides: `TELEGRAM_TOKEN`, `POSTGRES_DSN`, `REDIS_ADDR`.

---

## Bot creation API

`bot.New` accepts functional options - pass only what you need:

```go
b, err := bot.New(cfg.Telegram.Token,
    bot.WithDB(db),
    bot.WithRedis(redisClient),
    bot.WithI18n(locale),
    bot.WithAdminIDs(cfg.Telegram.AdminIDs),
    bot.WithLogger(logger),         // optional: default logger created automatically
    bot.WithWebhook(":8080"),       // optional: switch from polling to webhook
    bot.WithPollerTimeout(30*time.Second), // optional: default 10s
)
```

Start the bot (blocks until SIGINT/SIGTERM):

```go
b.RegisterHandlers()
b.Run()
```

---

## Adding handlers

All handlers live in `internal/bot/handler.go` → `RegisterHandlers()`:

```go
func (b *Bot) RegisterHandlers() {
    b.Use(middleware.Logger(b.logger))
    b.Use(middleware.RateLimit(30, time.Minute))

    b.Handle("/start", handlers.StartHandler(b.i18n))
    b.Handle("/help",  myHelpHandler)

    // admin-only
    b.Handle("/ban", banHandler, middleware.AdminOnly(b.adminIDs))
}
```

`b.Handle` returns `*Bot` for chaining. `b.Use` registers global middleware.

### Handler pattern

```go
// internal/bot/handlers/my_command.go
func MyHandler(i18n *i18n.I18n) telebot.HandlerFunc {
    return func(c telebot.Context) error {
        lang := c.Sender().LanguageCode
        if lang == "" {
            lang = "en"
        }
        return c.Send(i18n.Localize(lang, "my_key", nil))
    }
}
```

---

## Middleware (`pkg/middleware`)

```go
// Global - applies to all handlers
b.Use(middleware.Logger(logger))
b.Use(middleware.RateLimit(5, time.Minute))

// Per-handler
b.Handle("/admin", adminHandler, middleware.AdminOnly(adminIDs))

// Custom reject message
b.Use(middleware.RateLimit(3, time.Minute, func(c telebot.Context) error {
    return c.Send("Too many messages, please wait.")
}))
```

| Middleware                     | Description                                    |
| ------------------------------ | ---------------------------------------------- |
| `AdminOnly(ids, onReject?)`    | Allow only listed user IDs                     |
| `Logger(l)`                    | Log every update: user ID, name, text/callback |
| `RateLimit(n, per, onReject?)` | Per-user sliding window rate limiter           |

---

## FSM - conversation flows (`internal/fsm`)

```go
// in RegisterHandlers
f := fsm.New(fsm.NewRedisStorage(b.redis, fsm.WithTTL(3600))).
    Fallback(handlers.TextHandler(b.i18n)) // called when no state matches

b.Handle("/register", func(c telebot.Context) error {
    f.SetState(c, "ask_name")
    return c.Send("What's your name?")
})

b.Handle(telebot.OnText, f.Route(
    fsm.On("ask_name", func(c telebot.Context) error {
        f.SetStateData(c, "ask_age", c.Text()) // store name, move to next step
        return c.Send("How old are you?")
    }),
    fsm.On("ask_age", func(c telebot.Context) error {
        name, _ := f.GetData(c)               // retrieve name from previous step
        f.ClearState(c)
        return c.Send("Done, " + name + "!")
    }),
))
```

For testing use `fsm.NewMemoryStorage()` - no Redis required.

---

## Database

```go
// Auto-selects Postgres or SQLite based on config
db, err := database.FromConfig(ctx, cfg)

// Core operations
db.Exec(ctx, query, args...)
db.Query(ctx, query, args...)
db.QueryRow(ctx, query, args...)

// CRUD helpers
db.Insert(ctx, "users", map[string]any{"name": "Alice", "age": 30})
db.Update(ctx, "users", map[string]any{"age": 31}, "name = $1", "Alice")
db.Delete(ctx, "users", "id = $1", userID)
db.Select(ctx, "users", []string{"id", "name"}, "active = $1", true)
```

### Soft delete

For tables with a `deleted_at` column, type-assert to `database.SoftDeleteDatabase`:

```go
sdb := db.(database.SoftDeleteDatabase)
sdb.SoftDelete(ctx, "users", "id = $1", userID)  // sets deleted_at = now()
sdb.Restore(ctx,     "users", "id = $1", userID)  // clears deleted_at
sdb.HardDelete(ctx,  "users", "id = $1", userID)  // DELETE FROM ...
sdb.SelectDeleted(ctx, "users", cols, "id = $1", userID)
```

Both `PostgresDB` and `SQLiteDB` implement `SoftDeleteDatabase`.

---

## Payments - Telegram Stars

```go
product := handlers.StarProduct{
    Title:       "Premium",
    Description: "Unlock all features for 30 days",
    Payload:     "premium_1month", // returned in OnPayment - identify what was bought
    Stars:       100,
}

b.Handle("/buy",             handlers.SendInvoice(product))
b.Handle(telebot.OnCheckout, handlers.PreCheckout())         // auto-approve
b.Handle(telebot.OnPayment,  handlers.PaymentSuccess(b.db))  // records to DB + thanks user
```

Required table:

```sql
CREATE TABLE IF NOT EXISTS payments (
    id                 SERIAL PRIMARY KEY,
    user_id            BIGINT NOT NULL,
    telegram_charge_id TEXT   NOT NULL UNIQUE,
    payload            TEXT   NOT NULL,
    stars              INT    NOT NULL,
    created_at         TIMESTAMP DEFAULT NOW()
);
```

Add validation before `c.Accept()` in `PreCheckout()` if needed (inventory check, etc.).

---

## Redis

```go
// Key-value
redis.Set(ctx, "key", "value", ttlSeconds) // 0 = no expiry
value, _ := redis.Get(ctx, "key")
redis.Del(ctx, "key")
exists, _ := redis.Exists(ctx, "key")

// Atomic counters
n, _ := redis.Incr(ctx, "counter")

// Hash maps
redis.HSet(ctx, "user:42", "name", "Alice")
name, _ := redis.HGet(ctx, "user:42", "name")
all, _ := redis.HGetAll(ctx, "user:42")
```

Use `redis.NewMockClient()` in tests - no Redis server required.

---

## Testing

```bash
# All tests (Redis mocked - no server needed)
go test ./...

# Integration tests (requires Redis)
docker-compose up -d redis
go test ./...
```

---

## Docker

```bash
docker compose up -d        # start bot + Redis
docker compose logs -f bot  # tail logs
docker compose down         # stop
```

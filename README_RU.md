# tgbase - Go шаблон для Telegram-ботов

> 🇬🇧 [Read in English](README.md)

[![CI Status](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/ilyabrin/tgbase/branch/main/graph/badge.svg)](https://codecov.io/gh/ilyabrin/tgbase)
[![Go Report Card](https://goreportcard.com/badge/github.com/ilyabrin/tgbase)](https://goreportcard.com/report/github.com/ilyabrin/tgbase)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Удобный шаблон для разработки Telegram-ботов на Go. Клонируй, настрой, добавь хендлеры.

**Всё включено:** база данных (Postgres/SQLite), Redis, FSM, middleware, i18n, платежи (Telegram Stars), инлайн- и реплай-клавиатуры, пагинация, рассылка, deep links, инлайн-режим, Docker.

---

## Быстрый старт

```bash
git clone https://github.com/ilyabrin/tgbase && cd tgbase
go mod tidy
# укажи токен бота в config.yaml
go run cmd/app/main.go
```

---

## Структура проекта

```text
tgbase/
├── cmd/app/main.go               # точка входа - собирает всё воедино
├── config/
│   └── config.go                 # YAML-конфиг + переопределение через env
├── internal/
│   ├── bot/
│   │   ├── bot.go                # тип Bot, functional options, Run/Start
│   │   ├── broadcast.go          # Broadcast - рассылка списку пользователей
│   │   ├── handler.go            # RegisterHandlers - сюда добавляй хендлеры
│   │   └── handlers/
│   │       ├── inline.go         # инлайн-режим (OnQuery)
│   │       ├── payment.go        # платежи Telegram Stars
│   │       ├── registration.go   # демо многошагового FSM-флоу
│   │       ├── start.go          # /start
│   │       └── text.go           # фоллбэк на текстовые сообщения
│   ├── database/
│   │   ├── database.go           # интерфейсы Database и SoftDeleteDatabase
│   │   ├── postgres.go           # реализация PostgreSQL
│   │   └── sqlite.go             # реализация SQLite
│   ├── fsm/
│   │   ├── fsm.go                # FSM: New, Route, On, SetState, GetData…
│   │   └── storage.go            # интерфейс Storage, RedisStorage, MemoryStorage
│   ├── i18n/
│   │   └── i18n.go               # обёртка go-i18n, locales/*.yaml
│   └── redis/
│       ├── redis.go              # Client, Config, NewRedisClient, NewMockClient
│       └── real_redis.go         # реализация на go-redis/v9
├── pkg/
│   ├── deeplink/deeplink.go      # парсинг /start payload, генерация ссылок
│   ├── keyboard/keyboard.go      # построитель инлайн- и реплай-клавиатур
│   ├── logger/logger.go          # обёртка stdlib logger
│   ├── middleware/middleware.go  # AdminOnly, Logger, RateLimit, Recover
│   └── pagination/pagination.go  # навигация ← N/M → и помощник Page()
├── config.yaml
├── docker-compose.yml
└── Dockerfile
```

---

## Конфигурация (`config.yaml`)

```yaml
database:
  type: sqlite          # или "postgres"
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
  token: "ВАШ_ТОКЕН_БОТА"
  admin_ids:
    - 123456789          # ваш Telegram user ID
```

Переопределение через переменные окружения: `TELEGRAM_TOKEN`, `POSTGRES_DSN`, `REDIS_ADDR`.

---

## API создания бота

`bot.New` принимает функциональные опции - передавай только нужные:

```go
b, err := bot.New(cfg.Telegram.Token,
    bot.WithDB(db),
    bot.WithRedis(redisClient),
    bot.WithI18n(locale),
    bot.WithAdminIDs(cfg.Telegram.AdminIDs),
    bot.WithLogger(logger),
    bot.WithErrorHandler(func(err error, c telebot.Context) {
        if c != nil {
            c.Send("Что-то пошло не так, попробуй позже.")
        }
    }),
    bot.WithWebhook(":8080"),              // опционально: webhook вместо polling
    bot.WithPollerTimeout(30*time.Second), // опционально: по умолчанию 10s
)
```

Запуск бота (блокирует до SIGINT/SIGTERM):

```go
b.RegisterHandlers()
b.Run()
```

---

## Добавление хендлеров

Все хендлеры регистрируются в `internal/bot/handler.go` → `RegisterHandlers()`:

```go
func (b *Bot) RegisterHandlers() {
    b.Use(middleware.Logger(b.logger))
    b.Use(middleware.RateLimit(30, time.Minute))

    b.Handle("/start", handlers.StartHandler(b.i18n))
    b.Handle("/help",  myHelpHandler)

    // только для администраторов
    b.Handle("/ban", banHandler, middleware.AdminOnly(b.adminIDs))
}
```

`b.Handle` возвращает `*Bot` для чейнинга. `b.Use` регистрирует глобальный middleware.

### Паттерн хендлера

```go
// internal/bot/handlers/my_command.go
func MyHandler(i18n *i18n.I18n) telebot.HandlerFunc {
    return func(c telebot.Context) error {
        lang := c.Sender().LanguageCode
        if lang == "" {
            lang = "ru"
        }
        return c.Send(i18n.Localize(lang, "my_key", nil))
    }
}
```

---

## Middleware (`pkg/middleware`)

```go
// Глобально - применяется ко всем хендлерам
b.Use(middleware.Recover(logger))      // всегда первым
b.Use(middleware.Logger(logger))
b.Use(middleware.RateLimit(5, time.Minute))

// Для конкретного хендлера
b.Handle("/admin", adminHandler, middleware.AdminOnly(adminIDs))

// Кастомное сообщение при блокировке
b.Use(middleware.RateLimit(3, time.Minute, func(c telebot.Context) error {
    return c.Send("Слишком много сообщений, подождите немного.")
}))
```

| Middleware                     | Описание                                              |
| ------------------------------ | ----------------------------------------------------- |
| `Recover(l)`                   | Перехватывает панику в хендлере, логирует, возвр. err |
| `AdminOnly(ids, onReject?)`    | Пропускает только перечисленных пользователей         |
| `Logger(l)`                    | Логирует каждый апдейт: user ID, имя, текст/коллбэк   |
| `RateLimit(n, per, onReject?)` | Ограничение частоты сообщений на пользователя         |

---

## Рассылка (Broadcast)

Отправка сообщения списку пользователей с автоматическим rate limiting (по умолчанию 50 мс между отправками ≈ 20 сообщ/с, безопасно ниже лимита Telegram в 30 сообщ/с). Останавливается при отмене контекста.

```go
result := b.Broadcast(ctx, userIDs, "🔔 Текст объявления",
    bot.WithBroadcastDelay(50*time.Millisecond), // опционально, 50ms по умолчанию
    bot.WithBroadcastOnError(func(id int64, err error) {
        log.Printf("не удалось отправить пользователю %d: %v", id, err)
    }),
)
fmt.Printf("отправлено: %d, ошибок: %d\n", result.Sent, result.Failed)
```

---

## Пагинация (`pkg/pagination`)

Создание навигационных клавиатур ← N/M → для постраничного вывода данных из БД.

```go
import "tgbase/pkg/pagination"

pg := pagination.New(totalUsers, 10) // 10 записей на страницу

// В хендлере списка:
items := fetchPage(pg.Offset(page), pg.PageSize)
c.Send(renderList(items), pg.Keyboard(page, "users"))

// Регистрация коллбэков навигации:
b.Handle("\fusers_prev", usersPageHandler)
b.Handle("\fusers_next", usersPageHandler)

func usersPageHandler(c telebot.Context) error {
    page, _ := pagination.Page(c) // получить номер страницы из коллбэка
    items := fetchPage(pg.Offset(page), pg.PageSize)
    return c.Edit(renderList(items), pg.Keyboard(page, "users"))
}
```

| Функция                   | Описание                                    |
| ------------------------- | ------------------------------------------- |
| `New(total, pageSize)`    | Создать пейджер                             |
| `.Pages()`                | Общее количество страниц                    |
| `.Offset(page)`           | SQL OFFSET для страницы (нумерация с 1)     |
| `.Keyboard(page, prefix)` | Клавиатура ← N/M →; nil если страница одна  |
| `Page(c)`                 | Получить номер целевой страницы из коллбэка |

---

## Deep links (`pkg/deeplink`)

```go
import "tgbase/pkg/deeplink"

// /start - работает и как обычный /start, и со стартовым параметром
b.Handle("/start", func(c telebot.Context) error {
    payload := deeplink.Parse(c) // "" если /start без параметра
    if payload == "" {
        return c.Send("Добро пожаловать!")
    }
    return c.Send("Ты пришёл по ссылке: " + payload)
})

// Генерация реферальной ссылки
link := deeplink.Build("mybot", "ref_42")
// → "https://t.me/mybot?start=ref_42"
```

---

## Инлайн-режим

Позволяет пользователям набирать `@botname запрос` в любом чате и выбрать результат для отправки. Хендлер зарегистрирован в `RegisterHandlers`. Кастомизируй `InlineHandler()` в `internal/bot/handlers/inline.go`.

```go
func InlineHandler() telebot.HandlerFunc {
    return func(c telebot.Context) error {
        text := c.Query().Text
        results := telebot.Results{
            &telebot.ArticleResult{
                ResultBase:  telebot.ResultBase{ID: "r0"},
                Title:       text,
                Description: "Отправить в чат",
                Text:        text,
            },
        }
        return c.Answer(&telebot.QueryResponse{Results: results, CacheTime: 60})
    }
}
```

> **Внимание:** инлайн-режим нужно включить в [@BotFather](https://t.me/BotFather) → Bot Settings → Inline Mode.

---

## Клавиатуры (`pkg/keyboard`)

### Инлайн-клавиатура

Кнопки внутри сообщения.

```go
kb := keyboard.New().
    Row(keyboard.Btn("Купить", "btn_buy"), keyboard.Btn("Отмена", "btn_cancel")).
    Row(keyboard.URL("Сайт", "https://example.com")).
    Build()

c.Send("Выберите:", kb)

// Регистрация коллбэков
b.Handle("\fbtn_buy",    buyHandler)
b.Handle("\fbtn_cancel", cancelHandler)
```

| Функция                    | Описание                                         |
| -------------------------- | ------------------------------------------------ |
| `Btn(text, unique, data?)` | Кнопка с коллбэком; ключ хендлера: `"\f"+unique` |
| `URL(text, url)`           | Кнопка-ссылка                                    |
| `New()`                    | Создать построитель                              |
| `.Row(btns...)`            | Добавить ряд; поддерживает чейнинг               |
| `.Build()`                 | Вернуть `*telebot.ReplyMarkup`                   |

### Реплай-клавиатура

Постоянная панель кнопок над полем ввода.

```go
kb := keyboard.Reply().
    Row("Профиль", "Настройки").
    Row("Помощь").
    OneTime().
    Build()
c.Send("Выберите:", kb)

// Нажатие на реплай-кнопку отправляет её текст как обычное OnText-сообщение
b.Handle(telebot.OnText, func(c telebot.Context) error {
    switch c.Text() {
    case "Профиль":
        return c.Send("Ваш профиль...")
    case "Настройки":
        return c.Send("Меню настроек...")
    case "Помощь":
        return c.Send("Текст помощи...")
    }
    return nil
})

// Убрать клавиатуру
c.Send("Готово!", keyboard.Remove())
```

| Метод            | Описание                                          |
| ---------------- | ------------------------------------------------- |
| `Reply()`        | Создать построитель (resize включён по умолчанию) |
| `.Row(texts...)` | Добавить ряд текстовых кнопок                     |
| `.OneTime()`     | Скрыть после первого нажатия                      |
| `.Persistent()`  | Показывать между сообщениями                      |
| `.Placeholder()` | Подсказка в поле ввода                            |
| `.Build()`       | Вернуть `*telebot.ReplyMarkup`                    |
| `Remove()`       | Убрать текущую реплай-клавиатуру                  |

---

## FSM - диалоговые флоу (`internal/fsm`)

```go
f := fsm.New(fsm.NewRedisStorage(b.redis, fsm.WithTTL(3600))).
    Fallback(handlers.TextHandler(b.i18n))

b.Handle("/register", func(c telebot.Context) error {
    f.SetState(c, "ask_name")
    return c.Send("Как тебя зовут?")
})

b.Handle(telebot.OnText, f.Route(
    fsm.On("ask_name", func(c telebot.Context) error {
        f.SetStateData(c, "ask_age", c.Text()) // сохранить имя, перейти дальше
        return c.Send("Сколько тебе лет?")
    }),
    fsm.On("ask_age", func(c telebot.Context) error {
        name, _ := f.GetData(c)
        f.ClearState(c)
        return c.Send("Готово, " + name + "!")
    }),
))
```

В тестах используй `fsm.NewMemoryStorage()` - Redis не нужен.

---

## База данных

```go
// Автоматически выбирает Postgres или SQLite из конфига
db, err := database.FromConfig(ctx, cfg)

// Базовые операции
db.Exec(ctx, query, args...)
db.Query(ctx, query, args...)
db.QueryRow(ctx, query, args...)

// CRUD-помощники
db.Insert(ctx, "users", map[string]any{"name": "Алиса", "age": 30})
db.Update(ctx, "users", map[string]any{"age": 31}, "name = $1", "Алиса")
db.Delete(ctx, "users", "id = $1", userID)
db.Select(ctx, "users", []string{"id", "name"}, "active = $1", true)
```

### Мягкое удаление

Для таблиц с колонкой `deleted_at` - type-assert к `database.SoftDeleteDatabase`:

```go
sdb := db.(database.SoftDeleteDatabase)
sdb.SoftDelete(ctx, "users", "id = $1", userID)  // устанавливает deleted_at = now()
sdb.Restore(ctx,     "users", "id = $1", userID)  // сбрасывает deleted_at
sdb.HardDelete(ctx,  "users", "id = $1", userID)  // физическое удаление
sdb.SelectDeleted(ctx, "users", cols, "id = $1", userID)
```

Оба типа - `PostgresDB` и `SQLiteDB` - реализуют `SoftDeleteDatabase`.

---

## Платежи - Telegram Stars

```go
product := handlers.StarProduct{
    Title:       "Премиум",
    Description: "Доступ ко всем функциям на 30 дней",
    Payload:     "premium_1month",
    Stars:       100,
}

b.Handle("/buy",             handlers.SendInvoice(product))
b.Handle(telebot.OnCheckout, handlers.PreCheckout())
b.Handle(telebot.OnPayment,  handlers.PaymentSuccess(b.db))
```

Необходимая таблица:

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

---

## Redis

```go
redis.Set(ctx, "key", "value", ttlSeconds) // 0 = без срока действия
value, _ := redis.Get(ctx, "key")
redis.Del(ctx, "key")
exists, _ := redis.Exists(ctx, "key")

n, _ := redis.Incr(ctx, "counter")

redis.HSet(ctx, "user:42", "name", "Алиса")
name, _ := redis.HGet(ctx, "user:42", "name")
all, _ := redis.HGetAll(ctx, "user:42")
```

В тестах используй `redis.NewMockClient()` - Redis-сервер не нужен.

---

## Тестирование

```bash
# Все тесты (Redis замокан - сервер не нужен)
go test ./...

# С реальным Redis (интеграционные тесты)
docker-compose up -d redis
go test ./...
```

---

## Docker

```bash
docker compose up -d        # запустить бот + Redis
docker compose logs -f bot  # смотреть логи
docker compose down         # остановить
```

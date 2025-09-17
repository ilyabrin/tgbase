# Go Telegram Bot Template (tgbase)

[![CI Status](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml)
[![Redis Tests](https://github.com/ilyabrin/tgbase/actions/workflows/redis-tests.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/redis-tests.yml)

[![Codecov](https://codecov.io/gh/ilyabrin/tgbase/branch/main/graph/badge.svg)](https://codecov.io/gh/ilyabrin/tgbase)

[![Go Reference](https://pkg.go.dev/badge/tgbase)](https://pkg.go.dev/tgbase)

[![Go Report Card](https://goreportcard.com/badge/github.com/ilyabrin/tgbase)](https://goreportcard.com/report/github.com/ilyabrin/tgbase)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A clean, production-ready Go template for building Telegram bots with:

- **Database Support**: PostgreSQL and SQLite with automatic migrations
- **Redis Integration**: Optional Redis support for caching and session management
- **Configuration Management**: YAML-based config with environment variable overrides
- **Internationalization**: Multi-language support (EN/RU) with extensible localization
- **Telegram Bot Framework**: Built on `telebot` with structured command handlers
- **Docker Support**: Complete Docker and Docker Compose setup
- **Testing**: Comprehensive unit tests with mocking
- **Clean Architecture**: Organized code structure following Go best practices

## Установка зависимостей

```bash
go mod tidy
```

## Project Structure

```bash
tgbase/
├── cmd/app/              # Application entry point
│   └── main.go           # Main application file
├── config/               # Configuration management
│   ├── config.go         # Config loading and environment variables
│   └── config_test.go    # Configuration tests
├── internal/             # Private application code
│   ├── bot/              # Telegram bot implementation
│   │   ├── bot.go        # Bot initialization and lifecycle
│   │   ├── handler.go    # Handler registration
│   │   ├── bot_test.go   # Bot tests
│   │   └── handlers/     # Command handlers
│   │       ├── start.go  # /start command handler
│   │       ├── text.go   # Text message handler
│   │       └── redis2.go # Redis demo command
│   ├── database/         # Database layer
│   │   ├── database.go   # Database interface
│   │   ├── postgres.go   # PostgreSQL implementation
│   │   ├── sqlite.go     # SQLite implementation
│   │   └── *_test.go     # Database tests
│   ├── i18n/             # Internationalization
│   │   └── i18n.go       # Localization management
│   └── redis/            # Redis client wrapper
│       ├── redis.go      # Redis operations
│       └── redis_test.go # Redis tests
├── pkg/logger/           # Reusable logging package
│   └── logger.go         # Logger implementation
├── deploy/               # Deployment files
│   └── telegram-bot.service # Systemd service
├── docker-compose.yml    # Docker development environment
├── Dockerfile            # Container build configuration
└── config.yaml           # Application configuration
```

## Правила структуры проекта

1. **Тесты**:
   - Каждый пакет должен иметь тесты в файлах с суффиксом `_test.go`
   - Тесты должны быть написаны для всех публичных функций и методов
   - Используйте моки для внешних зависимостей

2. **Организация кода**:
   - Код должен быть организован по функциональности в соответствующих пакетах
   - Внутренние пакеты размещаются в директории `internal/`
   - Переиспользуемые пакеты размещаются в директории `pkg/`

3. **Зависимости**:
   - Используйте инъекцию зависимостей для тестируемости
   - Минимизируйте внешние зависимости

## Как использовать

1. **Инициализация проекта**:
   - Склонируйте структуру директорий.
   - Установите зависимости: `go mod tidy`.
   - Настройте `config.yaml` или `.env` с вашими параметрами.

2. **Добавление функционала**:
   - Расширяйте `internal/bot/handler.go` для обработки новых команд Telegram.
   - Добавляйте методы в `internal/database` для работы с данными.
   - Используйте Redis в `internal/redis` для кэширования или очередей.
   - Добавляйте новые языки в `i18n` для поддержки интернационализации в `locales/` (`en.yaml`, `ru.yaml`)

3. **Запуск**:

```bash
go run cmd/app/main.go
```

4. **Тестирование**:
    - Пишите тесты в соответствующих `*_test.go` файлах для каждого пакета.
    - Запускайте тесты с помощью `go test ./...`.

## Testing

### Running Tests

**Basic tests** (no external dependencies):

```bash
go test -short ./...
```

**All tests** (requires Redis):

```bash
# Start Redis first
docker-compose up -d redis

# Run all tests
go test ./...

# Clean up
docker-compose down
```

**Redis-specific tests**:

```bash
# Using project Redis container
docker-compose up -d redis
go test -v ./internal/redis/
docker-compose down

# Using standalone Redis
docker run -d --name redis-test -p 6379:6379 redis:7-alpine
go test -v ./internal/redis/
docker stop redis-test && docker rm redis-test

# Using test script
chmod +x test-redis.sh
./test-redis.sh
```

### GitHub Actions

The project includes comprehensive CI/CD with:

- **Main CI** (`ci.yml`): Runs on every push/PR with Redis services
- **Redis Tests** (`redis-tests.yml`): Matrix testing across Go/Redis versions
- **Release** (`release.yml`): Automated releases on version tags

The Redis integration tests run automatically on:

- Changes to `internal/redis/**`
- Changes to `go.mod`, `go.sum`
- Manual workflow dispatch

## Available Commands

The bot currently supports the following commands:

- `/start` - Welcome message with localization support
- `/redis2` - Redis demonstration with interactive toggle button (requires Redis)
- Text messages are echoed back with localization

## Adding New Commands

To add a new command to the bot:

1. **Create a handler** in `internal/bot/handlers/`:

   ```go
   // internal/bot/handlers/command_name.go
   package handlers

   import (
       "tgbase/internal/i18n"
       "gopkg.in/telebot.v3"
   )

   func CommandNameHandler(i18n *i18n.I18n) telebot.HandlerFunc {
       return func(c telebot.Context) error {
           lang := c.Sender().LanguageCode
           if lang == "" {
               lang = "en" // Fallback to English
           }
           message := i18n.Localize(lang, "command_response_key", nil)
           return c.Send(message)
       }
   }
   ```

2. **Register the command** in `internal/bot/handler.go`:

   ```go
   func (b *Bot) registerHandlers() {
       // Existing handlers
       b.bot.Handle("/start", handlers.StartHandler(b.i18n))
       b.bot.Handle(telebot.OnText, handlers.TextHandler(b.i18n))

       // Your new command
       b.bot.Handle("/command_name", handlers.CommandNameHandler(b.i18n))
   }
   ```

3. **Add tests** for the new command in `internal/bot/handlers/command_name_test.go`

### Рекомендации по командам

- **Именование команд**: Используйте понятные и короткие имена для команд
- **Интернационализация**: Всегда используйте i18n для текстовых сообщений
- **Обработка ошибок**: Корректно обрабатывайте ошибки и возвращайте понятные сообщения пользователю
- **Разделение ответственности**: Логика обработки команды должна быть отделена от бизнес-логики
- **Документация**: Документируйте новые команды в коде и README

### Типы обработчиков

Бот поддерживает различные типы обработчиков:

- **Команды**: `/start`, `/help`, и т.д.
- **Текстовые сообщения**: `telebot.OnText`
- **Кнопки**: Inline и Reply кнопки
- **Файлы**: Обработка загруженных пользователем файлов
- **Callback-запросы**: Обработка нажатий на inline-кнопки

## Docker Setup

### Local Development with Docker

1. Copy the sample config file and update it for Docker:

   ```bash
   cp config.yaml_sample config.docker.yaml
   ```

2. Edit `config.docker.yaml` and set your Telegram bot token:

   ```yaml
   telegram:
     token: "your-telegram-bot-token"
   ```

3. Start the services using Docker Compose:

   ```bash
   docker compose up -d
   ```

   This will start:
   - Redis server on port 6379
   - Your Telegram bot with hot-reload

4. Check the logs:

   ```bash
   docker compose logs -f
   ```

5. To stop the services:

   ```bash
   docker compose down
   ```

### Docker Commands

- Rebuild the bot after code changes:

  ```bash
  docker compose build bot
  ```

- Restart a specific service:

  ```bash
  docker compose restart bot
  ```

- View service logs:

  ```bash
  docker compose logs -f [service]  # bot or redis
  ```

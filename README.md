# Project Template for Go Telegram Bot (tgbase)

[![CI Status](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml/badge.svg)](https://github.com/ilyabrin/tgbase/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/ilyabrin/tgbase/branch/main/graph/badge.svg)](https://codecov.io/gh/ilyabrin/tgbase)

Go project template with:

- Configuration via YAML and environment variables
- Postgres/SQLite support
- Redis support (optional)
- Telegram bot integration using `telebot`
- Internationalization (i18n) with support for RU and EN languages
- CI/CD with GitHub Actions
- Systemd service for automatic bot restarts
- Unit tests for all components

## Установка зависимостей

```bash
go mod tidy
```

## Структура проекта

```
tgbase/
├── cmd/                  # Точки входа в приложение
│   └── app/              # Основное приложение
├── config/               # Конфигурация приложения
├── deploy/               # Файлы для деплоя (systemd, docker)
├── internal/             # Внутренний код приложения
│   ├── bot/              # Логика телеграм-бота
│   │   └── handlers/     # Обработчики команд бота
│   ├── database/         # Работа с базой данных
│   ├── i18n/             # Интернационализация
│   └── redis/            # Работа с Redis
└── pkg/                  # Переиспользуемые пакеты
    └── logger/           # Логирование
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

## Добавление новых команд бота

### Структура команд

Для добавления новой команды в бота следуйте этим шагам:

1. **Создайте обработчик команды** в директории `internal/bot/handlers/`:

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

2. **Зарегистрируйте команду** в файле `internal/bot/handler.go`:

   ```go
   func (b *Bot) registerHandlers() {
       // Существующие обработчики
       b.bot.Handle("/start", handlers.StartHandler(b.i18n))
       
       // Новая команда
       b.bot.Handle("/command_name", handlers.CommandNameHandler(b.i18n))
   }
   ```

3. **Добавьте тесты** для новой команды в `internal/bot/handlers/command_name_test.go`

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

# Project Template for Go Telegram Bot (tgbase)

Go project template with:

- Configuration via YAML and environment variables
- Postgres/SQLite support
- Redis support (optional)
- Telegram bot integration

## Установка зависимостей

```bash
go mod tidy
```

## Как использовать

1. **Инициализация проекта**:
   - Склонируйте структуру директорий.
   - Установите зависимости: `go mod tidy`.
   - Настройте `config.yaml` или `.env` с вашими параметрами.

2. **Добавление функционала**:
   - Расширяйте `internal/bot/handler.go` для обработки новых команд Telegram.
   - Добавляйте методы в `internal/database` для работы с данными.
   - Используйте Redis в `internal/redis` для кэширования или очередей.

3. **Запуск**:

```bash
go run cmd/app/main.go
```

4. **Тестирование**:
    - Пишите тесты в `internal/bot/handler_test.go` и `internal/database/database_test.go`.
    - Запускайте тесты с помощью `go test ./...`.

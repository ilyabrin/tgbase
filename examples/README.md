# Examples

This directory contains practical examples of how to implement common bot functionality using the tgbase framework.

## Available Examples

### 1. Basic Handlers (`basic_handler.go`)

- **BasicHandler**: Simple command that responds with a static message
- **EchoHandler**: Echoes user input with localization support

### 2. Database Integration (`database_example.go`)

- **UserStatsHandler**: Tracks user interaction statistics in the database
- **UserListHandler**: Queries and displays multiple database records

### 3. Redis Operations (`redis_example.go`)

- **CacheHandler**: Demonstrates caching expensive operations
- **SessionHandler**: Session management with user data persistence
- **CounterHandler**: Shared counter using Redis atomic operations

### 4. Interactive Keyboards (`inline_keyboard_example.go`)

- **MenuHandler**: Creates a main menu with multiple options
- **QuestionHandler**: Yes/No questions with callback handling
- **HandleCallbacks**: Universal callback handler for button interactions
- **DynamicMenuHandler**: Menu that changes based on user state/permissions

## How to Use These Examples

### 1. Copy the handler you need

```go
// Copy the function to your handlers directory
cp examples/basic_handler.go internal/bot/handlers/my_handler.go
```

### 2. Register the handler in your bot

```go
// In internal/bot/handler.go
func (b *Bot) registerHandlers() {
    // Existing handlers...

    // Add your example handler
    b.bot.Handle("/menu", examples.MenuHandler(b.i18n))
    b.bot.Handle("/stats", examples.UserStatsHandler(b.db, b.i18n))

    // Register callback handlers
    b.bot.Handle("btn_yes", examples.HandleCallbacks(b.redis, b.i18n))
    b.bot.Handle("btn_no", examples.HandleCallbacks(b.redis, b.i18n))
}
```

### 3. Customize for your needs

- Modify the response messages
- Add your own localization keys
- Extend the database schema
- Add new button types and callbacks

## Example Usage Patterns

### Database Operations

```go
// Create table
createTableQuery := `CREATE TABLE IF NOT EXISTS my_table (...)`
db.Exec(ctx, createTableQuery)

// Insert/Update with conflict resolution
insertQuery := `INSERT ... ON CONFLICT ... DO UPDATE SET ...`
db.Exec(ctx, insertQuery, params...)

// Query single row
row := db.QueryRow(ctx, "SELECT ... WHERE id = $1", userID)
row.Scan(&result)

// Query multiple rows
rows, _ := db.Query(ctx, "SELECT ... ORDER BY ...")
defer rows.Close()
for rows.Next() {
    rows.Scan(&data)
}
```

### Redis Operations

```go
// Basic key-value operations
redis.Set(ctx, "key", "value", expiration)
value, _ := redis.Get(ctx, "key")
redis.Del(ctx, "key")

// Hash operations for structured data
redis.HSet(ctx, "hash_key", "field", "value")
value, _ := redis.HGet(ctx, "hash_key", "field")
allData, _ := redis.HGetAll(ctx, "hash_key")

// Atomic counters
newValue, _ := redis.Incr(ctx, "counter_key")
```

### Inline Keyboards

```go
// Create keyboard
inlineKeys := &telebot.ReplyMarkup{}

// Create buttons
btn1 := inlineKeys.Data("Button Text", "callback_data")
btn2 := inlineKeys.Data("Another Button", "other_callback")

// Arrange in rows
inlineKeys.Inline(
    inlineKeys.Row(btn1, btn2),
    inlineKeys.Row(btn3),
)

// Send with message
c.Send("Choose an option:", inlineKeys)

// Handle callbacks
func handleCallback(c telebot.Context) error {
    defer c.Respond() // Always respond to remove loading state

    switch c.Callback().Data {
    case "callback_data":
        return c.Edit("Button 1 pressed!")
    case "other_callback":
        return c.Edit("Button 2 pressed!")
    }
    return nil
}
```

## Integration with Main Bot

To integrate these examples into your main bot, you need to:

1. Import the examples package in your handler registration
2. Add the handlers to your `registerHandlers()` function
3. Set up any required database tables
4. Configure Redis if using Redis-based examples
5. Add localization keys for any messages

Example integration:

```go
// internal/bot/handler.go
import "tgbase/examples"

func (b *Bot) registerHandlers() {
    // Basic commands
    b.bot.Handle("/start", handlers.StartHandler(b.i18n))

    // Example handlers
    b.bot.Handle("/menu", examples.MenuHandler(b.i18n))
    b.bot.Handle("/stats", examples.UserStatsHandler(b.db, b.i18n))

    // Redis examples (if Redis is available)
    if b.redis != nil {
        b.bot.Handle("/cache", examples.CacheHandler(b.redis, b.i18n))
        b.bot.Handle("/session", examples.SessionHandler(b.redis, b.i18n))
        b.bot.Handle("/counter", examples.CounterHandler(b.redis, b.i18n))
    }

    // Callback handlers
    b.bot.Handle(examples.BtnYes, examples.HandleCallbacks(b.redis, b.i18n))
    b.bot.Handle(examples.BtnNo, examples.HandleCallbacks(b.redis, b.i18n))
    b.bot.Handle(examples.BtnHelp, examples.HandleCallbacks(b.redis, b.i18n))
}

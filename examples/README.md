# Examples

Practical patterns for common bot functionality. Copy what you need into `internal/bot/handlers/` and register in `handler.go`.

## Available examples

### 1. Basic handlers (`basic_handler.go`)

- **BasicHandler** - static response to a command
- **EchoHandler** - echo user input with i18n

### 2. Database (`database_example.go`)

- **UserStatsHandler** - track user interactions (INSERT/UPDATE/SELECT)
- **UserListHandler** - paginated query results

### 3. Redis (`redis_example.go`)

- **CacheHandler** - cache expensive operations
- **SessionHandler** - per-user session data via hash maps
- **CounterHandler** - atomic shared counter

### 4. Inline keyboards (`inline_keyboard_example.go`)

- **MenuHandler** - main menu with multiple buttons
- **QuestionHandler** - yes/no with callback handling
- **HandleCallbacks** - universal callback dispatcher
- **DynamicMenuHandler** - state-aware menu

---

## How to use

### 1. Copy the handler

```bash
cp examples/basic_handler.go internal/bot/handlers/my_handler.go
```

### 2. Register in `handler.go`

```go
func (b *Bot) RegisterHandlers() {
    b.Handle("/menu",  myMenuHandler(b.i18n))
    b.Handle("/stats", myStatsHandler(b.db, b.i18n))
}
```

---

## Common patterns

### Database

```go
ctx := context.Background()

// Raw query
db.Exec(ctx, "CREATE TABLE IF NOT EXISTS my_table (...)")

// Insert with upsert
db.Exec(ctx, "INSERT INTO t (id, n) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET n = $2", id, n)

// Single row
row := db.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", userID)
row.Scan(&name)

// Multiple rows
rows, _ := db.Query(ctx, "SELECT id, name FROM users ORDER BY id")
defer rows.Close()
for rows.Next() {
    rows.Scan(&id, &name)
}
```

### Redis

```go
ctx := context.Background()

redis.Set(ctx, "key", "value", ttlSeconds) // 0 = no expiry
value, _ := redis.Get(ctx, "key")
redis.Del(ctx, "key")
exists, _ := redis.Exists(ctx, "key")

// Atomic counter
n, _ := redis.Incr(ctx, "counter")

// Hash map
redis.HSet(ctx, "user:42", "field", "value")
value, _ = redis.HGet(ctx, "user:42", "field")
all, _ := redis.HGetAll(ctx, "user:42")
```

### Inline keyboards

```go
markup := &telebot.ReplyMarkup{}
btn1 := markup.Data("Yes", "btn_yes")
btn2 := markup.Data("No",  "btn_no")
markup.Inline(markup.Row(btn1, btn2))

c.Send("Choose:", markup)

// Callback handler
func handleCallback(c telebot.Context) error {
    defer c.Respond() // remove loading spinner
    switch c.Callback().Data {
    case "btn_yes":
        return c.Edit("You chose Yes!")
    case "btn_no":
        return c.Edit("You chose No!")
    }
    return nil
}
```

### FSM (multi-step flows)

```go
f := fsm.New(fsm.NewRedisStorage(redisClient, fsm.WithTTL(3600))).
    Fallback(defaultTextHandler)

b.Handle("/start_flow", func(c telebot.Context) error {
    f.SetState(c, "step1")
    return c.Send("Step 1: enter something")
})

b.Handle(telebot.OnText, f.Route(
    fsm.On("step1", func(c telebot.Context) error {
        f.SetStateData(c, "step2", c.Text()) // save input, advance state
        return c.Send("Step 2: enter something else")
    }),
    fsm.On("step2", func(c telebot.Context) error {
        prev, _ := f.GetData(c)             // retrieve step1 input
        f.ClearState(c)
        return c.Send("Done! You entered: " + prev + " and " + c.Text())
    }),
))
```

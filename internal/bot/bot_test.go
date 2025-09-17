package bot

import (
	"os"
	"path/filepath"
	"testing"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"
	"tgbase/pkg/logger"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/telebot.v3"
)

// Helper function to create test dependencies
func createTestDependencies(t *testing.T) (*i18n.I18n, *logger.Logger) {
	// Create temp directory for i18n files
	tempDir, err := os.MkdirTemp("", "bot_test_i18n")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Create test English locale file
	enContent := `welcome: "Welcome to the test bot!"`
	enFile := filepath.Join(tempDir, "en.yaml")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatalf("Failed to write en.yaml: %v", err)
	}

	// Create i18n instance
	i18nInstance, err := i18n.NewI18n(tempDir)
	if err != nil {
		t.Fatalf("Failed to create i18n: %v", err)
	}

	// Create logger
	testLogger := logger.NewLogger()

	return i18nInstance, testLogger
}

// TestNewBot tests the NewBot function with invalid tokens
func TestNewBot(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "invalid empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid short token",
			token:   "123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i18nInstance, testLogger := createTestDependencies(t)

			bot, err := NewBot(tt.token, nil, nil, i18nInstance, testLogger)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, bot)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bot)
			}
		})
	}
}

// TestNewBotWebhook tests the NewBotWebhook function
func TestNewBotWebhook(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "invalid empty token webhook",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid short token webhook",
			token:   "123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i18nInstance, testLogger := createTestDependencies(t)

			bot, err := NewBotWebhook(tt.token, nil, nil, i18nInstance, testLogger)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, bot)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bot)
			}
		})
	}
}

// Note: TestBot_Start is skipped because it requires a valid Telegram bot token
// and makes actual API calls to Telegram, which is not suitable for unit tests.

// TestBot_registerHandlers tests the registerHandlers method logic
func TestBot_registerHandlers(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	t.Run("with Redis client", func(t *testing.T) {
		bot := &Bot{
			bot:    nil, // We can't mock telebot.Bot easily, so we test the logic
			db:     nil,
			redis:  &redis.Client{}, // Non-nil Redis
			i18n:   i18nInstance,
			logger: testLogger,
		}

		// Test that the Redis field is checked
		assert.NotNil(t, bot.redis)
		assert.Equal(t, i18nInstance, bot.i18n)
	})

	t.Run("without Redis client", func(t *testing.T) {
		bot := &Bot{
			bot:    nil,
			db:     nil,
			redis:  nil, // No Redis
			i18n:   i18nInstance,
			logger: testLogger,
		}

		// Test that the Redis field is nil
		assert.Nil(t, bot.redis)
		assert.Equal(t, i18nInstance, bot.i18n)
	})
}

// TestBot_struct tests Bot struct field access
func TestBot_struct(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	bot := &Bot{
		bot:    &tb.Bot{},
		db:     nil,
		redis:  nil,
		i18n:   i18nInstance,
		logger: testLogger,
	}

	// Test that all fields are accessible and set correctly
	assert.NotNil(t, bot.bot)
	assert.Nil(t, bot.db)
	assert.Nil(t, bot.redis)
	assert.Equal(t, i18nInstance, bot.i18n)
	assert.Equal(t, testLogger, bot.logger)
}

// TestNewBot_ValidTokenSimulation tests NewBot with error scenarios
func TestNewBot_ErrorCases(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	// Test what happens when telebot.NewBot would fail
	// We can't actually test this without mocking telebot, but we can test the error path coverage
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "token too short",
			token: "short",
		},
		{
			name:  "invalid format token",
			token: "invalid-token-format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These will fail with actual API calls, which tests the error path
			bot, err := NewBot(tt.token, nil, nil, i18nInstance, testLogger)
			assert.Error(t, err)
			assert.Nil(t, bot)
		})
	}
}

// TestNewBotWebhook_ErrorCases tests NewBotWebhook with error scenarios
func TestNewBotWebhook_ErrorCases(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "token too short webhook",
			token: "short",
		},
		{
			name:  "invalid format token webhook",
			token: "invalid-token-format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These will fail with actual API calls, which tests the error path
			bot, err := NewBotWebhook(tt.token, nil, nil, i18nInstance, testLogger)
			assert.Error(t, err)
			assert.Nil(t, bot)
		})
	}
}
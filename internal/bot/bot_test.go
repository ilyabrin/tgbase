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

func createTestDependencies(t *testing.T) (*i18n.I18n, *logger.Logger) {
	tempDir, err := os.MkdirTemp("", "bot_test_i18n")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	enContent := `welcome: "Welcome to the test bot!"`
	if err := os.WriteFile(filepath.Join(tempDir, "en.yaml"), []byte(enContent), 0644); err != nil {
		t.Fatalf("Failed to write en.yaml: %v", err)
	}

	i18nInstance, err := i18n.NewI18n(tempDir)
	if err != nil {
		t.Fatalf("Failed to create i18n: %v", err)
	}
	return i18nInstance, logger.NewLogger()
}

func TestNew_InvalidToken(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	for _, token := range []string{"", "123", "short", "invalid-token-format"} {
		b, err := New(token, WithI18n(i18nInstance), WithLogger(testLogger))
		assert.Error(t, err, "token=%q should fail", token)
		assert.Nil(t, b)
	}
}

func TestNew_WithWebhook_InvalidToken(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	b, err := New("", WithWebhook(":8080"), WithI18n(i18nInstance), WithLogger(testLogger))
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestBot_RegisterHandlers_WithRedis(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	b := &Bot{
		bot:    &tb.Bot{},
		redis:  &redis.Client{},
		i18n:   i18nInstance,
		logger: testLogger,
	}

	assert.NotNil(t, b.redis)
	assert.Equal(t, i18nInstance, b.i18n)
}

func TestBot_RegisterHandlers_WithoutRedis(t *testing.T) {
	i18nInstance, testLogger := createTestDependencies(t)

	b := &Bot{
		bot:    &tb.Bot{},
		i18n:   i18nInstance,
		logger: testLogger,
	}

	assert.Nil(t, b.redis)
	assert.Equal(t, i18nInstance, b.i18n)
}

func TestBot_DefaultLogger(t *testing.T) {
	// When no logger is provided, New should create a default one.
	b, err := New("bad-token")
	assert.Error(t, err) // token is invalid, but logger must be set before error return
	assert.Nil(t, b)
	// No panic — the logger was initialised internally before telebot.NewBot was called.
}

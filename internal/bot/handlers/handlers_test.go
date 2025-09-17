package handlers

import (
	"os"
	"path/filepath"
	"testing"
	"tgbase/internal/i18n"
	"tgbase/internal/redis"

	"github.com/stretchr/testify/assert"
)

// Helper function to create test i18n
func createTestI18n(t *testing.T) *i18n.I18n {
	tempDir, err := os.MkdirTemp("", "handlers_test_i18n")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	// Create test English locale file
	enContent := `welcome: "Welcome to the test bot!"
echo: "You said: {{.Text}}"`
	enFile := filepath.Join(tempDir, "en.yaml")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatalf("Failed to write en.yaml: %v", err)
	}

	// Create i18n instance
	i18nInstance, err := i18n.NewI18n(tempDir)
	if err != nil {
		t.Fatalf("Failed to create i18n: %v", err)
	}

	return i18nInstance
}

// TestStartHandler tests that StartHandler returns a valid handler function
func TestStartHandler(t *testing.T) {
	i18nInstance := createTestI18n(t)

	// Test that handler function is created successfully
	handler := StartHandler(i18nInstance)

	// Test that handler is not nil
	assert.NotNil(t, handler)

	// Test that it's a function type
	assert.NotNil(t, handler)
}

// TestTextHandler tests that TextHandler returns a valid handler function
func TestTextHandler(t *testing.T) {
	i18nInstance := createTestI18n(t)

	// Test that handler function is created successfully
	handler := TextHandler(i18nInstance)

	// Test that handler is not nil
	assert.NotNil(t, handler)

	// Test that it's a function type
	assert.NotNil(t, handler)
}

// TestRedis2Handler tests that Redis2Handler returns a function
func TestRedis2Handler(t *testing.T) {
	// Test that handler function is created successfully
	redisClient := &redis.Client{} // Empty client for testing
	handler := Redis2Handler(redisClient)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestHandleRedis2Button tests that HandleRedis2Button returns a function
func TestHandleRedis2Button(t *testing.T) {
	// Test that handler function is created successfully
	redisClient := &redis.Client{} // Empty client for testing
	handler := HandleRedis2Button(redisClient)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestHandlers_Constants tests that constants are properly defined
func TestHandlers_Constants(t *testing.T) {
	assert.Equal(t, "redis2", Redis2Key)
	assert.Equal(t, "btn_toggle", BtnToggle)
}

// TestStartHandler_WithNilI18n tests StartHandler with nil i18n
func TestStartHandler_WithNilI18n(t *testing.T) {
	// Test that handler function is created even with nil i18n
	handler := StartHandler(nil)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestTextHandler_WithNilI18n tests TextHandler with nil i18n
func TestTextHandler_WithNilI18n(t *testing.T) {
	// Test that handler function is created even with nil i18n
	handler := TextHandler(nil)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestRedis2Handler_WithNilRedis tests Redis2Handler with nil redis
func TestRedis2Handler_WithNilRedis(t *testing.T) {
	// Test that handler function is created even with nil redis
	handler := Redis2Handler(nil)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestHandleRedis2Button_WithNilRedis tests HandleRedis2Button with nil redis
func TestHandleRedis2Button_WithNilRedis(t *testing.T) {
	// Test that handler function is created even with nil redis
	handler := HandleRedis2Button(nil)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestStartHandler_WithEmptyI18n tests StartHandler with empty i18n
func TestStartHandler_WithEmptyI18n(t *testing.T) {
	// Create an empty i18n instance (no locale files)
	tempDir, err := os.MkdirTemp("", "handlers_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	i18nInstance, err := i18n.NewI18n(tempDir)
	if err != nil {
		t.Fatalf("Failed to create empty i18n: %v", err)
	}

	// Test that handler function is created successfully
	handler := StartHandler(i18nInstance)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}

// TestTextHandler_WithEmptyI18n tests TextHandler with empty i18n
func TestTextHandler_WithEmptyI18n(t *testing.T) {
	// Create an empty i18n instance (no locale files)
	tempDir, err := os.MkdirTemp("", "handlers_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	i18nInstance, err := i18n.NewI18n(tempDir)
	if err != nil {
		t.Fatalf("Failed to create empty i18n: %v", err)
	}

	// Test that handler function is created successfully
	handler := TextHandler(i18nInstance)

	// Test that handler is not nil
	assert.NotNil(t, handler)
}
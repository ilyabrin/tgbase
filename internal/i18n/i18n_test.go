package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewI18n_InvalidPath(t *testing.T) {
	// This won't error because filepath.Glob on non-existent directory just returns empty slice
	i18n, err := NewI18n("/nonexistent/path")
	if err != nil {
		t.Errorf("Unexpected error for nonexistent path: %v", err)
	}
	if i18n == nil {
		t.Error("Expected valid i18n instance, got nil")
	}
}

func TestNewI18n_EmptyPath(t *testing.T) {
	// This won't error either for the same reason
	i18n, err := NewI18n("")
	if err != nil {
		t.Errorf("Unexpected error for empty path: %v", err)
	}
	if i18n == nil {
		t.Error("Expected valid i18n instance, got nil")
	}
}

func TestI18n_Localize_WithoutFiles(t *testing.T) {
	// Create a temporary directory without any locale files
	tempDir, err := os.MkdirTemp("", "i18n_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	i18n, err := NewI18n(tempDir)
	if err != nil {
		t.Fatalf("NewI18n failed: %v", err)
	}

	// Test localization with no files - should return empty string
	result := i18n.Localize("en", "test_key", nil)
	if result != "" {
		t.Errorf("Expected empty string for missing localization, got '%s'", result)
	}
}

func TestI18n_Localize_WithTestFiles(t *testing.T) {
	// Create a temporary directory with test locale files
	tempDir, err := os.MkdirTemp("", "i18n_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test English locale file
	enContent := `welcome: "Welcome to the test bot!"
echo: "You said: {{.Text}}"
test_key: "Test value"`

	enFile := filepath.Join(tempDir, "en.yaml")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatalf("Failed to write en.yaml: %v", err)
	}

	// Create a test Russian locale file
	ruContent := `welcome: "Добро пожаловать в тестового бота!"
echo: "Вы сказали: {{.Text}}"
test_key: "Тестовое значение"`

	ruFile := filepath.Join(tempDir, "ru.yaml")
	if err := os.WriteFile(ruFile, []byte(ruContent), 0644); err != nil {
		t.Fatalf("Failed to write ru.yaml: %v", err)
	}

	i18n, err := NewI18n(tempDir)
	if err != nil {
		t.Fatalf("NewI18n failed: %v", err)
	}

	tests := []struct {
		name     string
		lang     string
		key      string
		data     map[string]any
		expected string
	}{
		{
			name:     "English simple key",
			lang:     "en",
			key:      "test_key",
			data:     nil,
			expected: "Test value",
		},
		{
			name:     "Russian simple key",
			lang:     "ru",
			key:      "test_key",
			data:     nil,
			expected: "Тестовое значение",
		},
		{
			name:     "English welcome",
			lang:     "en",
			key:      "welcome",
			data:     nil,
			expected: "Welcome to the test bot!",
		},
		{
			name:     "Russian welcome",
			lang:     "ru",
			key:      "welcome",
			data:     nil,
			expected: "Добро пожаловать в тестового бота!",
		},
		{
			name: "English echo with data",
			lang: "en",
			key:  "echo",
			data: map[string]any{
				"Text": "Hello World",
			},
			expected: "You said: Hello World",
		},
		{
			name: "Russian echo with data",
			lang: "ru",
			key:  "echo",
			data: map[string]any{
				"Text": "Привет мир",
			},
			expected: "Вы сказали: Привет мир",
		},
		{
			name:     "Nonexistent key",
			lang:     "en",
			key:      "nonexistent",
			data:     nil,
			expected: "",
		},
		{
			name:     "Nonexistent language falls back to English",
			lang:     "fr",
			key:      "test_key",
			data:     nil,
			expected: "Test value", // Falls back to English
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := i18n.Localize(tt.lang, tt.key, tt.data)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestI18n_Localize_EmptyData(t *testing.T) {
	// Create a temporary directory with test locale files
	tempDir, err := os.MkdirTemp("", "i18n_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with template
	enContent := `template_key: "Hello {{.Name}}!"`

	enFile := filepath.Join(tempDir, "en.yaml")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatalf("Failed to write en.yaml: %v", err)
	}

	i18n, err := NewI18n(tempDir)
	if err != nil {
		t.Fatalf("NewI18n failed: %v", err)
	}

	// Test with nil data
	result := i18n.Localize("en", "template_key", nil)
	if result != "Hello <no value>!" {
		t.Errorf("Expected template with <no value>, got '%s'", result)
	}

	// Test with empty data map
	result = i18n.Localize("en", "template_key", map[string]any{})
	if result != "Hello <no value>!" {
		t.Errorf("Expected template with <no value>, got '%s'", result)
	}
}
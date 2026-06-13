package handlers

import (
	"os"
	"path/filepath"
	"testing"
	"tgbase/internal/fsm"
	"tgbase/internal/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

// fakeCtx is a minimal telebot.Context for handler dispatch tests.
type fakeCtx struct {
	telebot.Context
	userID   int64
	langCode string
	text     string
	sent     []interface{}
	store    map[string]interface{}
}

func (f *fakeCtx) Sender() *telebot.User {
	return &telebot.User{ID: f.userID, LanguageCode: f.langCode}
}
func (f *fakeCtx) Text() string { return f.text }
func (f *fakeCtx) Send(what interface{}, _ ...interface{}) error {
	f.sent = append(f.sent, what)
	return nil
}
func (f *fakeCtx) Set(key string, val interface{}) {
	if f.store == nil {
		f.store = map[string]interface{}{}
	}
	f.store[key] = val
}
func (f *fakeCtx) Get(key string) interface{} {
	if f.store == nil {
		return nil
	}
	return f.store[key]
}

func newCtx(userID int64, lang, text string) *fakeCtx {
	return &fakeCtx{userID: userID, langCode: lang, text: text}
}

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

// --- Dispatch tests ---

func TestStartHandler_Dispatch(t *testing.T) {
	h := StartHandler(createTestI18n(t))

	t.Run("known lang", func(t *testing.T) {
		c := newCtx(1, "en", "")
		require.NoError(t, h(c))
		assert.Len(t, c.sent, 1)
		assert.Equal(t, "Welcome to the test bot!", c.sent[0])
	})

	t.Run("empty lang falls back to en", func(t *testing.T) {
		c := newCtx(1, "", "")
		require.NoError(t, h(c))
		assert.Len(t, c.sent, 1)
		assert.Equal(t, "Welcome to the test bot!", c.sent[0])
	})
}

func TestTextHandler_Dispatch(t *testing.T) {
	h := TextHandler(createTestI18n(t))

	c := newCtx(1, "en", "hello")
	require.NoError(t, h(c))
	assert.Len(t, c.sent, 1)
}

func TestRegisterStart_SetsStateAndReplies(t *testing.T) {
	f := fsm.New(fsm.NewMemoryStorage())
	h := RegisterStart(f)
	c := newCtx(42, "", "")
	require.NoError(t, h(c))

	state, err := f.GetState(c)
	require.NoError(t, err)
	assert.Equal(t, StateAskName, state)
	assert.Len(t, c.sent, 1)
}

func TestRegisterAskName_ValidName(t *testing.T) {
	f := fsm.New(fsm.NewMemoryStorage())
	h := RegisterAskName(f)
	c := newCtx(42, "", "Alice")
	require.NoError(t, h(c))

	state, _ := f.GetState(c)
	assert.Equal(t, StateAskAge, state)
	data, _ := f.GetData(c)
	assert.Equal(t, "Alice", data)
	assert.Len(t, c.sent, 1)
}

func TestRegisterAskName_EmptyName(t *testing.T) {
	f := fsm.New(fsm.NewMemoryStorage())
	h := RegisterAskName(f)
	c := newCtx(42, "", "   ")
	require.NoError(t, h(c))

	state, _ := f.GetState(c)
	assert.Equal(t, fsm.None, state, "state should not advance on empty name")
	assert.Len(t, c.sent, 1)
}

func TestRegisterAskAge_ValidAge(t *testing.T) {
	f := fsm.New(fsm.NewMemoryStorage())
	setup := newCtx(42, "", "")
	f.SetStateData(setup, StateAskAge, "Bob")

	h := RegisterAskAge(f)
	c := newCtx(42, "", "25")
	require.NoError(t, h(c))

	state, _ := f.GetState(c)
	assert.Equal(t, fsm.None, state, "state should be cleared after completion")
	require.Len(t, c.sent, 1)
	assert.Contains(t, c.sent[0], "Bob")
}

func TestRegisterAskAge_InvalidAge(t *testing.T) {
	f := fsm.New(fsm.NewMemoryStorage())
	h := RegisterAskAge(f)

	for _, bad := range []string{"abc", "0", "200", "-5"} {
		c := newCtx(42, "", bad)
		require.NoError(t, h(c), "input=%q", bad)
		assert.Len(t, c.sent, 1, "input=%q", bad)
	}
}

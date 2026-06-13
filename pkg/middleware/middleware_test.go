package middleware

import (
	"testing"
	"time"

	"tgbase/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

// mockContext implements enough of telebot.Context to drive middleware.
type mockContext struct {
	telebot.Context
	senderID int64
	text     string
	called   bool
}

func (m *mockContext) Sender() *telebot.User       { return &telebot.User{ID: m.senderID} }
func (m *mockContext) Callback() *telebot.Callback { return nil }
func (m *mockContext) Message() *telebot.Message {
	return &telebot.Message{Text: m.text}
}
func (m *mockContext) Send(what interface{}, opts ...interface{}) error { return nil }

func nextHandler(called *bool) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		*called = true
		return nil
	}
}

// --- AdminOnly ---

func TestAdminOnly_AllowsAdmin(t *testing.T) {
	mw := AdminOnly([]int64{100, 200})
	called := false
	handler := mw(nextHandler(&called))

	err := handler(&mockContext{senderID: 100})
	require.NoError(t, err)
	assert.True(t, called, "admin should reach next handler")
}

func TestAdminOnly_BlocksNonAdmin(t *testing.T) {
	mw := AdminOnly([]int64{100})
	called := false
	handler := mw(nextHandler(&called))

	err := handler(&mockContext{senderID: 999})
	require.NoError(t, err)
	assert.False(t, called, "non-admin should be blocked")
}

func TestAdminOnly_OnRejectCalled(t *testing.T) {
	rejected := false
	mw := AdminOnly([]int64{100}, func(c telebot.Context) error {
		rejected = true
		return nil
	})
	handler := mw(nextHandler(new(bool)))

	err := handler(&mockContext{senderID: 999})
	require.NoError(t, err)
	assert.True(t, rejected, "onReject should be called for non-admin")
}

func TestAdminOnly_EmptyList(t *testing.T) {
	mw := AdminOnly(nil)
	called := false
	handler := mw(nextHandler(&called))

	err := handler(&mockContext{senderID: 1})
	require.NoError(t, err)
	assert.False(t, called)
}

// --- Logger ---

func TestLogger_PassesThrough(t *testing.T) {
	l := logger.NewLogger()
	mw := Logger(l)
	called := false
	handler := mw(nextHandler(&called))

	err := handler(&mockContext{senderID: 1, text: "hello"})
	require.NoError(t, err)
	assert.True(t, called, "Logger must call next handler")
}

// --- RateLimit ---

func TestRateLimit_AllowsUnderLimit(t *testing.T) {
	mw := RateLimit(3, time.Minute)
	handler := mw(func(c telebot.Context) error { return nil })

	ctx := &mockContext{senderID: 1}
	for i := 0; i < 3; i++ {
		require.NoError(t, handler(ctx), "message %d should be allowed", i+1)
	}
}

func TestRateLimit_BlocksOverLimit(t *testing.T) {
	blocked := 0
	mw := RateLimit(2, time.Minute, func(c telebot.Context) error {
		blocked++
		return nil
	})
	handler := mw(func(c telebot.Context) error { return nil })

	ctx := &mockContext{senderID: 42}
	for i := 0; i < 5; i++ {
		require.NoError(t, handler(ctx))
	}
	assert.Equal(t, 3, blocked, "messages 3-5 should be rate-limited")
}

func TestRateLimit_ResetsAfterWindow(t *testing.T) {
	mw := RateLimit(1, 50*time.Millisecond)
	calls := 0
	handler := mw(func(c telebot.Context) error { calls++; return nil })

	ctx := &mockContext{senderID: 7}

	require.NoError(t, handler(ctx)) // allowed
	require.NoError(t, handler(ctx)) // blocked

	time.Sleep(60 * time.Millisecond)

	require.NoError(t, handler(ctx)) // new window - allowed
	assert.Equal(t, 2, calls)
}

func TestRateLimit_IndependentUsers(t *testing.T) {
	mw := RateLimit(1, time.Minute)
	calls := map[int64]int{}
	handler := mw(func(c telebot.Context) error {
		calls[c.Sender().ID]++
		return nil
	})

	handler(&mockContext{senderID: 1})
	handler(&mockContext{senderID: 2})
	handler(&mockContext{senderID: 1}) // blocked for user 1
	handler(&mockContext{senderID: 2}) // blocked for user 2

	assert.Equal(t, 1, calls[1])
	assert.Equal(t, 1, calls[2])
}

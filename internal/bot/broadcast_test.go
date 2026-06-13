package bot

import (
	"context"
	"errors"
	"testing"
	"time"

	"tgbase/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

// mockSender records Send calls and can simulate errors.
type mockSender struct {
	sent   []int64
	failID int64 // if non-zero, return an error for this user ID
}

func (m *mockSender) Send(to telebot.Recipient, _ interface{}, _ ...interface{}) (*telebot.Message, error) {
	id := to.(telebot.ChatID)
	m.sent = append(m.sent, int64(id))
	if m.failID != 0 && int64(id) == m.failID {
		return nil, errors.New("blocked")
	}
	return &telebot.Message{}, nil
}

func cfg0() *broadcastCfg { return &broadcastCfg{delay: 0} }

func TestBroadcast_AllSucceed(t *testing.T) {
	s := &mockSender{}
	res := runBroadcast(context.Background(), s, logger.NewLogger(), []int64{1, 2, 3}, "hi", cfg0())

	assert.Equal(t, 3, res.Sent)
	assert.Equal(t, 0, res.Failed)
	assert.Empty(t, res.Errors)
	assert.Equal(t, []int64{1, 2, 3}, s.sent)
}

func TestBroadcast_PartialFailure(t *testing.T) {
	s := &mockSender{failID: 2}
	res := runBroadcast(context.Background(), s, logger.NewLogger(), []int64{1, 2, 3}, "hi", cfg0())

	assert.Equal(t, 2, res.Sent)
	assert.Equal(t, 1, res.Failed)
	require.Len(t, res.Errors, 1)
	assert.Equal(t, int64(2), res.Errors[0].UserID)
	assert.Error(t, res.Errors[0].Err)
}

func TestBroadcast_OnErrorCallback(t *testing.T) {
	s := &mockSender{failID: 2}
	var cbID int64
	cfg := &broadcastCfg{delay: 0, onError: func(id int64, _ error) { cbID = id }}

	runBroadcast(context.Background(), s, logger.NewLogger(), []int64{1, 2}, "hi", cfg)
	assert.Equal(t, int64(2), cbID)
}

func TestBroadcast_ContextCancellation(t *testing.T) {
	s := &mockSender{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before start

	res := runBroadcast(ctx, s, logger.NewLogger(), []int64{1, 2, 3}, "hi", cfg0())
	assert.Equal(t, 0, res.Sent, "no sends after context is already cancelled")
}

func TestBroadcast_EmptyList(t *testing.T) {
	s := &mockSender{}
	res := runBroadcast(context.Background(), s, logger.NewLogger(), nil, "hi", cfg0())
	assert.Equal(t, 0, res.Sent)
	assert.Equal(t, 0, res.Failed)
}

func TestWithBroadcastDelay(t *testing.T) {
	cfg := &broadcastCfg{}
	WithBroadcastDelay(200 * time.Millisecond)(cfg)
	assert.Equal(t, 200*time.Millisecond, cfg.delay)
}

func TestWithBroadcastOnError(t *testing.T) {
	cfg := &broadcastCfg{}
	called := false
	WithBroadcastOnError(func(int64, error) { called = true })(cfg)
	cfg.onError(0, nil)
	assert.True(t, called)
}

package fsm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestMemoryStorage(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()

	t.Run("default state is None", func(t *testing.T) {
		state, err := s.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, None, state)
	})

	t.Run("set and get", func(t *testing.T) {
		require.NoError(t, s.Set(ctx, 1, "ask_name"))
		state, err := s.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, State("ask_name"), state)
	})

	t.Run("clear returns to None", func(t *testing.T) {
		s.Set(ctx, 2, "ask_age")
		require.NoError(t, s.Clear(ctx, 2))
		state, err := s.Get(ctx, 2)
		require.NoError(t, err)
		assert.Equal(t, None, state)
	})

	t.Run("independent users", func(t *testing.T) {
		s.Set(ctx, 10, "step_a")
		s.Set(ctx, 20, "step_b")
		a, _ := s.Get(ctx, 10)
		b, _ := s.Get(ctx, 20)
		assert.Equal(t, State("step_a"), a)
		assert.Equal(t, State("step_b"), b)
	})
}

func noop(telebot.Context) error { return nil }

func TestFSM_Route_buildsTable(t *testing.T) {
	f := New(NewMemoryStorage())
	h := f.Route(On("step1", noop), On("step2", noop))
	assert.NotNil(t, h)
}

func TestFSM_Fallback(t *testing.T) {
	f := New(NewMemoryStorage())
	assert.Nil(t, f.fallback)
	f.Fallback(noop)
	assert.NotNil(t, f.fallback)
}

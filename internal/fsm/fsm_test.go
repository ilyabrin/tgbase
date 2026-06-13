package fsm

import (
	"context"
	"testing"

	"tgbase/internal/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

// fakeContext is a minimal telebot.Context for FSM tests.
type fakeContext struct {
	telebot.Context
	userID int64
	text   string
}

func (f *fakeContext) Sender() *telebot.User                       { return &telebot.User{ID: f.userID} }
func (f *fakeContext) Text() string                                { return f.text }
func (f *fakeContext) Send(_ interface{}, _ ...interface{}) error  { return nil }

func ctx42() *fakeContext { return &fakeContext{userID: 42} }

// --- MemoryStorage ---

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

// --- splitRaw ---

func TestSplitRaw(t *testing.T) {
	t.Run("no data", func(t *testing.T) {
		state, data := splitRaw("ask_name")
		assert.Equal(t, "ask_name", state)
		assert.Equal(t, "", data)
	})

	t.Run("with data", func(t *testing.T) {
		state, data := splitRaw("ask_age" + sep + "John")
		assert.Equal(t, "ask_age", state)
		assert.Equal(t, "John", data)
	})

	t.Run("empty raw", func(t *testing.T) {
		state, data := splitRaw("")
		assert.Equal(t, None, state)
		assert.Equal(t, "", data)
	})

	t.Run("data may contain sep - only first sep is the delimiter", func(t *testing.T) {
		state, data := splitRaw("step" + sep + "a" + sep + "b")
		assert.Equal(t, "step", state)
		assert.Equal(t, "a"+sep+"b", data)
	})
}

// --- SetState / GetState ---

func TestFSM_SetAndGetState(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	state, err := f.GetState(c)
	require.NoError(t, err)
	assert.Equal(t, None, state)

	require.NoError(t, f.SetState(c, "step1"))
	state, err = f.GetState(c)
	require.NoError(t, err)
	assert.Equal(t, State("step1"), state)
}

func TestFSM_ClearState(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	f.SetState(c, "step1")
	require.NoError(t, f.ClearState(c))

	state, err := f.GetState(c)
	require.NoError(t, err)
	assert.Equal(t, None, state)
}

// --- SetStateData / GetData ---

func TestFSM_SetStateData_GetData(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	require.NoError(t, f.SetStateData(c, "ask_age", "Alice"))

	state, err := f.GetState(c)
	require.NoError(t, err)
	assert.Equal(t, State("ask_age"), state)

	data, err := f.GetData(c)
	require.NoError(t, err)
	assert.Equal(t, "Alice", data)
}

func TestFSM_GetData_NoData(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	f.SetState(c, "ask_name") // no data
	data, err := f.GetData(c)
	require.NoError(t, err)
	assert.Equal(t, "", data)
}

func TestFSM_ClearState_RemovesData(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	f.SetStateData(c, "ask_age", "Bob")
	f.ClearState(c)

	data, _ := f.GetData(c)
	assert.Equal(t, "", data)
}

// --- Route dispatch ---

func TestFSM_Route_MatchesState(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()

	called := ""
	handler := f.Route(
		On("step1", func(tc telebot.Context) error { called = "step1"; return nil }),
		On("step2", func(tc telebot.Context) error { called = "step2"; return nil }),
	)

	f.SetState(c, "step1")
	require.NoError(t, handler(c))
	assert.Equal(t, "step1", called)

	f.SetState(c, "step2")
	require.NoError(t, handler(c))
	assert.Equal(t, "step2", called)
}

func TestFSM_Route_NoMatchCallsFallback(t *testing.T) {
	fallbackCalled := false
	f := New(NewMemoryStorage()).
		Fallback(func(telebot.Context) error { fallbackCalled = true; return nil })

	c := ctx42()
	// state is None — no step matches
	handler := f.Route(On("step1", func(telebot.Context) error { return nil }))
	require.NoError(t, handler(c))
	assert.True(t, fallbackCalled)
}

func TestFSM_Route_NoMatchNoFallback_Silent(t *testing.T) {
	f := New(NewMemoryStorage())
	c := ctx42()
	handler := f.Route(On("step1", func(telebot.Context) error { return nil }))
	// state is None, no fallback — should return nil silently
	assert.NoError(t, handler(c))
}

func TestFSM_Route_IndependentUsers(t *testing.T) {
	f := New(NewMemoryStorage())
	c1 := &fakeContext{userID: 1}
	c2 := &fakeContext{userID: 2}

	calls := map[int64]bool{}
	handler := f.Route(
		On("active", func(tc telebot.Context) error {
			calls[tc.Sender().ID] = true
			return nil
		}),
	)

	f.SetState(c1, "active")
	// c2 stays at None

	handler(c1)
	handler(c2)

	assert.True(t, calls[1], "user 1 should hit the handler")
	assert.False(t, calls[2], "user 2 should not hit the handler")
}

// --- RedisStorage ---

func TestRedisStorage(t *testing.T) {
	ctx := context.Background()

	t.Run("default state is None", func(t *testing.T) {
		s := NewRedisStorage(redis.NewMockClient())
		val, err := s.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, None, val)
	})

	t.Run("set and get", func(t *testing.T) {
		s := NewRedisStorage(redis.NewMockClient())
		require.NoError(t, s.Set(ctx, 1, "step1"))
		val, err := s.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, State("step1"), val)
	})

	t.Run("clear", func(t *testing.T) {
		s := NewRedisStorage(redis.NewMockClient())
		s.Set(ctx, 1, "step1")
		require.NoError(t, s.Clear(ctx, 1))
		val, err := s.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, None, val)
	})

	t.Run("WithPrefix changes key namespace", func(t *testing.T) {
		client := redis.NewMockClient()
		s1 := NewRedisStorage(client, WithPrefix("ns1"))
		s2 := NewRedisStorage(client, WithPrefix("ns2"))
		s1.Set(ctx, 42, "a")
		val, _ := s2.Get(ctx, 42)
		assert.Equal(t, None, val, "different prefixes should not share state")
	})

	t.Run("WithTTL sets expiry field", func(t *testing.T) {
		s := NewRedisStorage(redis.NewMockClient(), WithTTL(3600))
		assert.Equal(t, int64(3600), s.ttl)
	})

	t.Run("independent users", func(t *testing.T) {
		s := NewRedisStorage(redis.NewMockClient())
		s.Set(ctx, 10, "a")
		s.Set(ctx, 20, "b")
		a, _ := s.Get(ctx, 10)
		b, _ := s.Get(ctx, 20)
		assert.Equal(t, State("a"), a)
		assert.Equal(t, State("b"), b)
	})
}

// --- Fallback setter ---

func noop(telebot.Context) error { return nil }

func TestFSM_Fallback(t *testing.T) {
	f := New(NewMemoryStorage())
	assert.Nil(t, f.fallback)
	f.Fallback(noop)
	assert.NotNil(t, f.fallback)
}

// --- StateData round-trip ---

func TestFSM_StateData_RoundTrip(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()

	encoded := "ask_age" + sep + "Alice"
	require.NoError(t, s.Set(ctx, 42, encoded))

	state, data := splitRaw(encoded)
	assert.Equal(t, "ask_age", state)
	assert.Equal(t, "Alice", data)

	require.NoError(t, s.Clear(ctx, 42))
	raw, _ := s.Get(ctx, 42)
	st, d := splitRaw(raw)
	assert.Equal(t, None, st)
	assert.Equal(t, "", d)
}

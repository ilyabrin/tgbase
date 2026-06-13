package fsm

import (
	"context"
	"fmt"
	"sync"

	"tgbase/internal/redis"
)

// Storage persists per-user conversation states.
type Storage interface {
	Get(ctx context.Context, userID int64) (State, error)
	Set(ctx context.Context, userID int64, state State) error
	Clear(ctx context.Context, userID int64) error
}

// --- RedisStorage ---

// RedisStorage stores FSM state in Redis.
type RedisStorage struct {
	client *redis.Client
	prefix string
	ttl    int64 // seconds; 0 = no expiry
}

type StorageOption func(*RedisStorage)

func WithPrefix(prefix string) StorageOption {
	return func(s *RedisStorage) { s.prefix = prefix }
}

// WithTTL sets state expiry in seconds (e.g. 3600 for 1 hour).
func WithTTL(seconds int64) StorageOption {
	return func(s *RedisStorage) { s.ttl = seconds }
}

// NewRedisStorage creates a Redis-backed state storage.
// Default key prefix is "fsm", no TTL.
func NewRedisStorage(client *redis.Client, opts ...StorageOption) *RedisStorage {
	s := &RedisStorage{client: client, prefix: "fsm"}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *RedisStorage) key(userID int64) string {
	return fmt.Sprintf("%s:%d", s.prefix, userID)
}

func (s *RedisStorage) Get(ctx context.Context, userID int64) (State, error) {
	val, err := s.client.Get(ctx, s.key(userID))
	if err != nil {
		return None, err
	}
	return State(val), nil
}

func (s *RedisStorage) Set(ctx context.Context, userID int64, state State) error {
	return s.client.Set(ctx, s.key(userID), string(state), s.ttl)
}

func (s *RedisStorage) Clear(ctx context.Context, userID int64) error {
	return s.client.Del(ctx, s.key(userID))
}

// --- MemoryStorage ---

// MemoryStorage stores FSM state in memory. Intended for testing.
type MemoryStorage struct {
	mu     sync.RWMutex
	states map[int64]State
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{states: make(map[int64]State)}
}

func (m *MemoryStorage) Get(_ context.Context, userID int64) (State, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.states[userID], nil
}

func (m *MemoryStorage) Set(_ context.Context, userID int64, state State) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[userID] = state
	return nil
}

func (m *MemoryStorage) Clear(_ context.Context, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, userID)
	return nil
}

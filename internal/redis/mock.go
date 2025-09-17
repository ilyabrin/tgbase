package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// MockRedisClient provides an in-memory implementation of Redis for testing
type MockRedisClient struct {
	mu     sync.RWMutex
	data   map[string]string
	hashes map[string]map[string]string
	expiry map[string]time.Time
	closed bool
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:   make(map[string]string),
		hashes: make(map[string]map[string]string),
		expiry: make(map[string]time.Time),
	}
}

func (m *MockRedisClient) checkExpiry(key string) {
	if expTime, exists := m.expiry[key]; exists {
		if time.Now().After(expTime) {
			delete(m.data, key)
			delete(m.hashes, key)
			delete(m.expiry, key)
		}
	}
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value any, expiration int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("redis client is closed")
	}

	m.data[key] = fmt.Sprintf("%v", value)

	if expiration > 0 {
		m.expiry[key] = time.Now().Add(time.Duration(expiration) * time.Second)
	} else {
		delete(m.expiry, key)
	}

	return nil
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return "", fmt.Errorf("redis client is closed")
	}

	m.checkExpiry(key)

	value, exists := m.data[key]
	if !exists {
		return "", nil
	}

	return value, nil
}

func (m *MockRedisClient) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("redis client is closed")
	}

	delete(m.data, key)
	delete(m.hashes, key)
	delete(m.expiry, key)

	return nil
}

func (m *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return false, fmt.Errorf("redis client is closed")
	}

	m.checkExpiry(key)

	_, exists := m.data[key]
	if !exists {
		_, exists = m.hashes[key]
	}

	return exists, nil
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return 0, fmt.Errorf("redis client is closed")
	}

	m.checkExpiry(key)

	value, exists := m.data[key]
	var intVal int64 = 0

	if exists {
		var err error
		intVal, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("value is not an integer")
		}
	}

	intVal++
	m.data[key] = strconv.FormatInt(intVal, 10)

	return intVal, nil
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return 0, fmt.Errorf("redis client is closed")
	}

	m.checkExpiry(key)

	value, exists := m.data[key]
	var intVal int64 = 0

	if exists {
		var err error
		intVal, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("value is not an integer")
		}
	}

	intVal--
	m.data[key] = strconv.FormatInt(intVal, 10)

	return intVal, nil
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, field string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("redis client is closed")
	}

	if m.hashes[key] == nil {
		m.hashes[key] = make(map[string]string)
	}

	m.hashes[key][field] = fmt.Sprintf("%v", value)

	return nil
}

func (m *MockRedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return "", fmt.Errorf("redis client is closed")
	}

	hash, exists := m.hashes[key]
	if !exists {
		return "", nil
	}

	value, exists := hash[field]
	if !exists {
		return "", nil
	}

	return value, nil
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, field string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("redis client is closed")
	}

	hash, exists := m.hashes[key]
	if !exists {
		return nil
	}

	delete(hash, field)

	if len(hash) == 0 {
		delete(m.hashes, key)
	}

	return nil
}

func (m *MockRedisClient) HExists(ctx context.Context, key string, field string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return false, fmt.Errorf("redis client is closed")
	}

	hash, exists := m.hashes[key]
	if !exists {
		return false, nil
	}

	_, exists = hash[field]
	return exists, nil
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, fmt.Errorf("redis client is closed")
	}

	hash, exists := m.hashes[key]
	if !exists {
		return make(map[string]string), nil
	}

	// Return a copy to avoid concurrent modification
	result := make(map[string]string)
	for k, v := range hash {
		result[k] = v
	}

	return result, nil
}

func (m *MockRedisClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	m.data = nil
	m.hashes = nil
	m.expiry = nil

	return nil
}

// Ping simulates a Redis ping command
func (m *MockRedisClient) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return fmt.Errorf("redis client is closed")
	}

	return nil
}
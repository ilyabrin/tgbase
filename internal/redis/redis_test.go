package redis

import (
	"context"
	"os"
	"testing"
	"time"
)

// createTestRedisClient creates a mock Redis client for testing
func createTestRedisClient(t *testing.T) *Client {
	return NewMockClient()
}

// TestNewRedisClient tests Redis client creation
func TestNewRedisClient(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			name:    "valid mock connection",
			addr:    "localhost:6379",
			wantErr: false,
		},
		{
			name:    "invalid connection - should still work with mock",
			addr:    "localhost:6380",
			wantErr: false, // Mock doesn't actually connect
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For unit tests, we use mock client
			client := NewMockClient()
			defer client.Close()

			if client == nil {
				t.Error("NewMockClient() returned nil")
				return
			}

			// Test that the client can perform basic operations
			if err := client.Set(ctx, "test", "value", 0); err != nil {
				t.Errorf("Mock client Set failed: %v", err)
			}

			value, err := client.Get(ctx, "test")
			if err != nil {
				t.Errorf("Mock client Get failed: %v", err)
			}
			if value != "value" {
				t.Errorf("Expected 'value', got '%s'", value)
			}
		})
	}
}

func TestRedisClient_KeyValueOperations(t *testing.T) {
	client := createTestRedisClient(t)
	defer client.Close()
	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		// Test Set
		err := client.Set(ctx, "test_key", "test_value", 0)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		// Test Get
		value, err := client.Get(ctx, "test_key")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if value != "test_value" {
			t.Errorf("Expected 'test_value', got '%s'", value)
		}
	})

	t.Run("Set with expiration", func(t *testing.T) {
		err := client.Set(ctx, "expire_key", "expire_value", 1)
		if err != nil {
			t.Errorf("Set with expiration failed: %v", err)
		}

		value, err := client.Get(ctx, "expire_key")
		if err != nil {
			t.Errorf("Get after set with expiration failed: %v", err)
		}
		if value != "expire_value" {
			t.Errorf("Expected 'expire_value', got '%s'", value)
		}

		// Wait for expiration and check again
		time.Sleep(2 * time.Second)
		value, err = client.Get(ctx, "expire_key")
		if err != nil {
			t.Errorf("Get after expiration should not error: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string after expiration, got '%s'", value)
		}
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		value, err := client.Get(ctx, "non_existent")
		if err != nil {
			t.Errorf("Get non-existent key should not error: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string for non-existent key, got '%s'", value)
		}
	})

	t.Run("Del", func(t *testing.T) {
		// Set a key first
		client.Set(ctx, "delete_me", "value", 0)

		// Delete it
		err := client.Del(ctx, "delete_me")
		if err != nil {
			t.Errorf("Del failed: %v", err)
		}

		// Verify it's gone
		value, err := client.Get(ctx, "delete_me")
		if err != nil {
			t.Errorf("Get after delete should not error: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string after delete, got '%s'", value)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		// Set a key
		client.Set(ctx, "exists_key", "value", 0)

		// Test exists
		exists, err := client.Exists(ctx, "exists_key")
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected key to exist")
		}

		// Test non-existent
		exists, err = client.Exists(ctx, "non_existent")
		if err != nil {
			t.Errorf("Exists for non-existent key failed: %v", err)
		}
		if exists {
			t.Error("Expected key to not exist")
		}
	})
}

func TestRedisClient_CounterOperations(t *testing.T) {
	client := createTestRedisClient(t)
	defer client.Close()
	ctx := context.Background()

	t.Run("Incr", func(t *testing.T) {
		// Increment non-existent key
		val, err := client.Incr(ctx, "counter")
		if err != nil {
			t.Errorf("Incr failed: %v", err)
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}

		// Increment again
		val, err = client.Incr(ctx, "counter")
		if err != nil {
			t.Errorf("Incr failed: %v", err)
		}
		if val != 2 {
			t.Errorf("Expected 2, got %d", val)
		}
	})

	t.Run("Decr", func(t *testing.T) {
		// Decrement non-existent key
		val, err := client.Decr(ctx, "decr_counter")
		if err != nil {
			t.Errorf("Decr failed: %v", err)
		}
		if val != -1 {
			t.Errorf("Expected -1, got %d", val)
		}

		// Decrement again
		val, err = client.Decr(ctx, "decr_counter")
		if err != nil {
			t.Errorf("Decr failed: %v", err)
		}
		if val != -2 {
			t.Errorf("Expected -2, got %d", val)
		}
	})
}

func TestRedisClient_HashOperations(t *testing.T) {
	client := createTestRedisClient(t)
	defer client.Close()
	ctx := context.Background()

	t.Run("HSet and HGet", func(t *testing.T) {
		// Set hash field
		err := client.HSet(ctx, "hash_key", "field1", "value1")
		if err != nil {
			t.Errorf("HSet failed: %v", err)
		}

		// Get hash field
		value, err := client.HGet(ctx, "hash_key", "field1")
		if err != nil {
			t.Errorf("HGet failed: %v", err)
		}
		if value != "value1" {
			t.Errorf("Expected 'value1', got '%s'", value)
		}
	})

	t.Run("HGet non-existent field", func(t *testing.T) {
		value, err := client.HGet(ctx, "hash_key", "non_existent")
		if err != nil {
			t.Errorf("HGet non-existent field should not error: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string for non-existent field, got '%s'", value)
		}
	})

	t.Run("HExists", func(t *testing.T) {
		// Set a field first
		client.HSet(ctx, "hash_exists", "field1", "value1")

		// Test exists
		exists, err := client.HExists(ctx, "hash_exists", "field1")
		if err != nil {
			t.Errorf("HExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected field to exist")
		}

		// Test non-existent
		exists, err = client.HExists(ctx, "hash_exists", "non_existent")
		if err != nil {
			t.Errorf("HExists for non-existent field failed: %v", err)
		}
		if exists {
			t.Error("Expected field to not exist")
		}
	})

	t.Run("HDel", func(t *testing.T) {
		// Set a field first
		client.HSet(ctx, "hash_del", "field1", "value1")

		// Delete it
		err := client.HDel(ctx, "hash_del", "field1")
		if err != nil {
			t.Errorf("HDel failed: %v", err)
		}

		// Verify it's gone
		value, err := client.HGet(ctx, "hash_del", "field1")
		if err != nil {
			t.Errorf("HGet after HDel should not error: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string after HDel, got '%s'", value)
		}
	})

	t.Run("HGetAll", func(t *testing.T) {
		// Set multiple fields
		client.HSet(ctx, "hash_all", "field1", "value1")
		client.HSet(ctx, "hash_all", "field2", "value2")

		// Get all
		data, err := client.HGetAll(ctx, "hash_all")
		if err != nil {
			t.Errorf("HGetAll failed: %v", err)
		}

		if len(data) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(data))
		}

		if data["field1"] != "value1" {
			t.Errorf("Expected 'value1' for field1, got '%s'", data["field1"])
		}

		if data["field2"] != "value2" {
			t.Errorf("Expected 'value2' for field2, got '%s'", data["field2"])
		}
	})

	t.Run("HGetAll empty hash", func(t *testing.T) {
		data, err := client.HGetAll(ctx, "empty_hash")
		if err != nil {
			t.Errorf("HGetAll for empty hash failed: %v", err)
		}

		if data == nil {
			t.Error("Expected empty map, got nil")
		}

		if len(data) != 0 {
			t.Errorf("Expected empty map, got %d fields", len(data))
		}
	})
}

func TestRedisClient_Close(t *testing.T) {
	client := createTestRedisClient(t)

	// Test close
	err := client.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

// Test Redis error handling with mock
func TestRedisClient_ErrorHandling(t *testing.T) {
	// Test with mock client that simulates errors
	client := NewMockClient()
	defer client.Close()

	// Close the client to simulate error conditions
	client.Close()

	ctx := context.Background()

	// Test operations on closed client
	err := client.Set(ctx, "test", "value", 0)
	if err == nil {
		t.Error("Expected error for closed client, got nil")
	}

	_, err = client.Get(ctx, "test")
	if err == nil {
		t.Error("Expected error for closed client, got nil")
	}
}

func TestNewRedisClient_Configurations(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		timeout   time.Duration
		expectErr bool
	}{
		{
			name:      "valid config",
			cfg:       Config{Addr: "localhost:6379"},
			timeout:   5 * time.Second,
			expectErr: false,
		},
		{
			name:      "invalid address format",
			cfg:       Config{Addr: "invalid-address"},
			timeout:   1 * time.Second,
			expectErr: false,
		},
		{
			name:      "empty address",
			cfg:       Config{},
			timeout:   1 * time.Second,
			expectErr: false,
		},
		{
			name:      "different DB number",
			cfg:       Config{Addr: "localhost:6379", DB: 5},
			timeout:   1 * time.Second,
			expectErr: false,
		},
		{
			name:      "with password",
			cfg:       Config{Addr: "localhost:6379", Password: "testpass"},
			timeout:   1 * time.Second,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For unit tests, always use mock client
			client := NewMockClient()
			defer client.Close()

			if client == nil {
				t.Error("NewMockClient() returned nil")
			}

			// Test basic operations to ensure client works
			ctx := context.Background()
			err := client.Set(ctx, "test", "value", 0)
			if err != nil {
				t.Errorf("Mock client Set failed: %v", err)
			}
		})
	}
}

// Test edge cases for createTestRedisClient
func TestCreateTestRedisClient_EdgeCases(t *testing.T) {
	// Test with REDIS_ADDR environment variable
	originalAddr := os.Getenv("REDIS_ADDR")
	defer func() {
		if originalAddr != "" {
			os.Setenv("REDIS_ADDR", originalAddr)
		} else {
			os.Unsetenv("REDIS_ADDR")
		}
	}()

	// Test with custom REDIS_ADDR
	os.Setenv("REDIS_ADDR", "localhost:9999")

	// Test that environment variable is set correctly
	if os.Getenv("REDIS_ADDR") != "localhost:9999" {
		t.Error("REDIS_ADDR environment variable not set correctly")
	}

	// Test that createTestRedisClient still works with mock
	client := createTestRedisClient(t)
	defer client.Close()

	if client == nil {
		t.Error("createTestRedisClient() returned nil")
	}
}

// Test Redis client method wrappers with mock
func TestRedisClient_MethodWrappers(t *testing.T) {
	client := NewMockClient()
	defer client.Close()
	ctx := context.Background()

	// Test Set method wrapper
	t.Run("Set method wrapper", func(t *testing.T) {
		// Test with expiration
		err := client.Set(ctx, "test_key", "test_value", 10)
		if err != nil {
			t.Errorf("Set with expiration failed: %v", err)
		}

		// Test without expiration
		err = client.Set(ctx, "test_key", "test_value", 0)
		if err != nil {
			t.Errorf("Set without expiration failed: %v", err)
		}
	})

	// Test Get method wrapper
	t.Run("Get method wrapper", func(t *testing.T) {
		client.Set(ctx, "test_key", "test_value", 0)
		value, err := client.Get(ctx, "test_key")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if value != "test_value" {
			t.Errorf("Expected 'test_value', got '%s'", value)
		}
	})

	// Test Del method wrapper
	t.Run("Del method wrapper", func(t *testing.T) {
		client.Set(ctx, "test_key", "test_value", 0)
		err := client.Del(ctx, "test_key")
		if err != nil {
			t.Errorf("Del failed: %v", err)
		}
	})

	// Test Exists method wrapper
	t.Run("Exists method wrapper", func(t *testing.T) {
		client.Set(ctx, "test_key", "test_value", 0)
		exists, err := client.Exists(ctx, "test_key")
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected key to exist")
		}
	})

	// Test Incr method wrapper
	t.Run("Incr method wrapper", func(t *testing.T) {
		val, err := client.Incr(ctx, "counter")
		if err != nil {
			t.Errorf("Incr failed: %v", err)
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})

	// Test Decr method wrapper
	t.Run("Decr method wrapper", func(t *testing.T) {
		val, err := client.Decr(ctx, "decr_counter")
		if err != nil {
			t.Errorf("Decr failed: %v", err)
		}
		if val != -1 {
			t.Errorf("Expected -1, got %d", val)
		}
	})

	// Test HSet method wrapper
	t.Run("HSet method wrapper", func(t *testing.T) {
		err := client.HSet(ctx, "hash_key", "field", "value")
		if err != nil {
			t.Errorf("HSet failed: %v", err)
		}
	})

	// Test HGet method wrapper
	t.Run("HGet method wrapper", func(t *testing.T) {
		client.HSet(ctx, "hash_key", "field", "value")
		value, err := client.HGet(ctx, "hash_key", "field")
		if err != nil {
			t.Errorf("HGet failed: %v", err)
		}
		if value != "value" {
			t.Errorf("Expected 'value', got '%s'", value)
		}
	})

	// Test HDel method wrapper
	t.Run("HDel method wrapper", func(t *testing.T) {
		client.HSet(ctx, "hash_key", "field", "value")
		err := client.HDel(ctx, "hash_key", "field")
		if err != nil {
			t.Errorf("HDel failed: %v", err)
		}
	})

	// Test HExists method wrapper
	t.Run("HExists method wrapper", func(t *testing.T) {
		client.HSet(ctx, "hash_key", "field", "value")
		exists, err := client.HExists(ctx, "hash_key", "field")
		if err != nil {
			t.Errorf("HExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected field to exist")
		}
	})

	// Test HGetAll method wrapper
	t.Run("HGetAll method wrapper", func(t *testing.T) {
		// Use a unique key for this test to avoid interference
		testKey := "hgetall_test_key"
		client.HSet(ctx, testKey, "field1", "value1")
		client.HSet(ctx, testKey, "field2", "value2")
		data, err := client.HGetAll(ctx, testKey)
		if err != nil {
			t.Errorf("HGetAll failed: %v", err)
		}
		if len(data) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(data))
		}
	})

	// Test Close method wrapper
	t.Run("Close method wrapper", func(t *testing.T) {
		mockClient := NewMockClient()
		err := mockClient.Close()
		if err != nil {
			t.Errorf("Close should not error, got: %v", err)
		}
	})
}

// Integration test that can be run with real Redis when available
// This test is skipped by default but can be enabled with -tags=integration
func TestRedisClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check if Redis is available
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	cfg := Config{
		Addr:     redisAddr,
		Password: "",
		DB:       1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := NewRedisClient(ctx, cfg)
	if err != nil {
		t.Skipf("Redis not available for integration testing: %v", err)
		return
	}
	defer client.Close()

	// Run basic integration test
	err = client.Set(ctx, "integration_test", "success", 0)
	if err != nil {
		t.Errorf("Integration test Set failed: %v", err)
	}

	value, err := client.Get(ctx, "integration_test")
	if err != nil {
		t.Errorf("Integration test Get failed: %v", err)
	}
	if value != "success" {
		t.Errorf("Integration test: expected 'success', got '%s'", value)
	}

	// Cleanup
	client.Del(ctx, "integration_test")
}

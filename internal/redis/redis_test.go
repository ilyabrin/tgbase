package redis

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tests := []struct {
		name string
		cfg  struct {
			Addr     string
			Password string
			DB       int
		}
		wantErr bool
	}{
		{
			name: "valid connection",
			cfg: struct {
				Addr     string
				Password string
				DB       int
			}{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
			wantErr: false,
		},
		{
			name: "invalid connection",
			cfg: struct {
				Addr     string
				Password string
				DB       int
			}{
				Addr:     "localhost:6380", // Wrong port
				Password: "",
				DB:       0,
			},
			wantErr: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRedisClient(ctx, tt.cfg)
			if err == nil {
				defer client.Close()
				// Additional validation when client is successfully created
				if client.client == nil {
					t.Error("NewRedisClient() returned nil redis client")
					return
				}
				// Verify connection is working
				if err := client.client.Ping(ctx).Err(); err != nil {
					t.Errorf("Redis connection test failed: %v", err)
				}
			}
		})
	}
}

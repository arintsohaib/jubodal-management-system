package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisManager handles session storage in Redis
type RedisManager struct {
	client *redis.Client
}

// NewRedisManager creates a new Redis manager
func NewRedisManager(redisURL, password string) (*RedisManager, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	if password != "" {
		opts.Password = password
	}

	client := redis.NewClient(opts)

	// Verify connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisManager{client: client}, nil
}

// SetSession stores a user session
func (m *RedisManager) SetSession(ctx context.Context, userID, tokenID string, expiry time.Duration) error {
	key := fmt.Sprintf("session:%s:%s", userID, tokenID)
	return m.client.Set(ctx, key, "active", expiry).Err()
}

// ValidateSession checks if a session exists and is active
func (m *RedisManager) ValidateSession(ctx context.Context, userID, tokenID string) (bool, error) {
	key := fmt.Sprintf("session:%s:%s", userID, tokenID)
	val, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "active", nil
}

// InvalidateSession removes a session
func (m *RedisManager) InvalidateSession(ctx context.Context, userID, tokenID string) error {
	key := fmt.Sprintf("session:%s:%s", userID, tokenID)
	return m.client.Del(ctx, key).Err()
}

// InvalidateAllUserSessions removes all sessions for a user
func (m *RedisManager) InvalidateAllUserSessions(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("session:%s:*", userID)
	
	iter := m.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := m.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	
	return iter.Err()
}

// IsRevoked checks if a token has been explicitly revoked
func (m *RedisManager) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("revoked:%s", tokenID)
	exists, err := m.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// RevokeToken adds a token to the revocation list
func (m *RedisManager) RevokeToken(ctx context.Context, tokenID string, expiry time.Duration) error {
	key := fmt.Sprintf("revoked:%s", tokenID)
	return m.client.Set(ctx, key, "revoked", expiry).Err()
}

// Client returns the underlying redis client
func (m *RedisManager) Client() *redis.Client {
	return m.client
}

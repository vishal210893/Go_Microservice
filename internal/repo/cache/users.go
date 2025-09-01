// Package cache provides Redis-based caching implementations for user entities.
package cache

import (
	"Go-Microservice/internal/repo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
	// ttl is the time-to-live for cached user entries
	ttl time.Duration
}

// NewUserStore creates a new UserStore with the specified Redis client and TTL.
//
// Parameters:
//   - rdb: Redis client instance
//   - ttl: Time-to-live for cached entries
//
// Returns:
//   - *UserStore: Configured user store instance
func NewUserStore(rdb *redis.Client, ttl time.Duration) *UserStore {
	return &UserStore{
		rdb: rdb,
		ttl: ttl,
	}
}

// Get retrieves a user from cache by ID.
// Returns nil if the user is not found in cache.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: Unique identifier for the user
//
// Returns:
//   - *repo.User: User data if found, nil if not found
//   - error: Error if operation fails
func (s *UserStore) Get(ctx context.Context, userID int64) (*repo.User, error) {
	cacheKey := s.buildCacheKey(userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	if data == "" {
		return nil, nil
	}

	var user repo.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &user, nil
}

// Set stores a user in cache with the configured TTL.
//
// Parameters:
//   - ctx: Context for the operation
//   - user: User data to cache
//
// Returns:
//   - error: Error if operation fails
func (s *UserStore) Set(ctx context.Context, user *repo.User) error {
	return s.SetWithTTL(ctx, user, s.ttl)
}

// SetWithTTL stores a user in cache with a custom TTL.
//
// Parameters:
//   - ctx: Context for the operation
//   - user: User data to cache
//   - ttl: Custom time-to-live for this entry
//
// Returns:
//   - error: Error if operation fails
func (s *UserStore) SetWithTTL(ctx context.Context, user *repo.User, ttl time.Duration) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	cacheKey := s.buildCacheKey(user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	if err := s.rdb.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set user in cache: %w", err)
	}

	return nil
}

// Delete removes a user from cache by ID.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: Unique identifier for the user to delete
func (s *UserStore) Delete(ctx context.Context, userID int64) {
	cacheKey := s.buildCacheKey(userID)
	s.rdb.Del(ctx, cacheKey).Err()
}

// Exists checks if a user exists in cache by ID.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: Unique identifier for the user
//
// Returns:
//   - bool: True if user exists in cache, false otherwise
//   - error: Error if operation fails
func (s *UserStore) Exists(ctx context.Context, userID int64) (bool, error) {
	cacheKey := s.buildCacheKey(userID)

	count, err := s.rdb.Exists(ctx, cacheKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check user existence in cache: %w", err)
	}

	return count > 0, nil
}

// buildCacheKey constructs a standardized cache key for the given user ID.
//
// Parameters:
//   - userID: User identifier
//
// Returns:
//   - string: Formatted cache key
func (s *UserStore) buildCacheKey(userID int64) string {
	return fmt.Sprintf("user-%d", userID)
}

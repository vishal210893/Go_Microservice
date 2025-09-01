// Package cache provides Redis-based storage implementations for caching
// repository entities such as users, posts, and comments.
package cache

import (
	"Go-Microservice/internal/repo"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Get(ctx context.Context, userID int64) (*repo.User, error)
	Set(ctx context.Context, user *repo.User) error
	Delete(ctx context.Context, userID int64)
}

type Storage struct {
	Users UserCache
	rdb   *redis.Client
}

// NewRedisStorage creates and returns a new Redis-based cache storage instance.
// It initializes cache stores for different entity types with the provided Redis client.
//
// Parameters:
//   - rdb: Redis client instance for cache operations
//
// Returns:
//   - Storage: Configured cache storage with initialized entity stores
//
// Example:
//
//	redisClient, err := NewRedisClient("localhost:6379", "", 0)
//	if err != nil {
//	    log.Fatal("Failed to create Redis client:", err)
//	}
//
//	cacheStorage := NewRedisStorage(redisClient)
//	user, err := cacheStorage.Users.Get(ctx, userID)
func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: NewUserStore(rdb, time.Hour),
	}
}

// Close closes the underlying Redis client connection.
// This method should be called when the storage is no longer needed
// to properly release resources.
//
// Returns:
//   - error: Error if closing the connection fails
func (s *Storage) Close() error {
	return s.rdb.Close()
}

// Ping tests the Redis connection to ensure it's still active.
// This can be used for health checks and connection validation.
//
// Parameters:
//   - ctx: Context for the ping operation
//
// Returns:
//   - error: Error if the ping fails or connection is down
func (s *Storage) Ping(ctx context.Context) error {
	return s.rdb.Ping(ctx).Err()
}

// FlushAll removes all keys from all Redis databases.
// This method should be used with caution, typically only in testing environments.
//
// Parameters:
//   - ctx: Context for the flush operation
//
// Returns:
//   - error: Error if the flush operation fails
func (s *Storage) FlushAll(ctx context.Context) error {
	return s.rdb.FlushAll(ctx).Err()
}

// Package cache provides Redis client configuration and connection utilities
// for caching operations in the application.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Default Redis client configuration constants
const (
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout = 3 * time.Second
	DefaultWriteTimeout = 3 * time.Second
	DefaultPoolSize = 10
	DefaultMinIdleConns = 5
)

// RedisConfig holds the configuration parameters for Redis client.
type RedisConfig struct {
	Addr string
	Password string
	DB int
	// PoolSize is the maximum number of socket connections
	PoolSize int
	// MinIdleConns is the minimum number of idle connections
	MinIdleConns int
	// ConnectTimeout is the timeout for establishing connections
	ConnectTimeout time.Duration
	// ReadTimeout is the timeout for socket reads
	ReadTimeout time.Duration
	// WriteTimeout is the timeout for socket writes
	WriteTimeout time.Duration
}

// NewRedisClient creates and returns a new Redis client with the specified configuration.
// It establishes a connection to the Redis server and validates the connection.
//
// Parameters:
//   - addr: Redis server address in format "host:port"
//   - pw: Redis server password (empty string if no authentication required)
//   - db: Redis database number (typically 0-15)
//
// Returns:
//   - *redis.Client: Configured Redis client instance
//   - error: Connection or configuration error if any
//
// Example:
//   client, err := NewRedisClient("localhost:6379", "", 0)
//   if err != nil {
//       log.Fatal("Failed to connect to Redis:", err)
//   }
//   defer client.Close()
func NewRedisClient(addr, pw string, db int) (*redis.Client, error) {
	config := RedisConfig{
		Addr:           addr,
		Password:       pw,
		DB:             db,
		PoolSize:       DefaultPoolSize,
		MinIdleConns:   DefaultMinIdleConns,
		ConnectTimeout: DefaultConnectTimeout,
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,
	}

	return NewRedisClientWithConfig(config)
}

// NewRedisClientWithConfig creates a Redis client with custom configuration.
// This function provides more control over Redis client settings including
// connection pooling, timeouts, and other advanced options.
//
// Parameters:
//   - config: RedisConfig struct containing all Redis client settings
//
// Returns:
//   - *redis.Client: Configured Redis client instance
//   - error: Connection or configuration error if any
func NewRedisClientWithConfig(config RedisConfig) (*redis.Client, error) {
	if config.Addr == "" {
		return nil, fmt.Errorf("redis address cannot be empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.ConnectTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", config.Addr, err)
	}

	return client, nil
}

package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFixedWindowRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisFixedWindowLimiter(client *redis.Client, limit int, window time.Duration) *RedisFixedWindowRateLimiter {
	return &RedisFixedWindowRateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisFixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:%s", ip)

	// Check current count (equivalent to RLock + read)
	count, err := rl.client.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		// Fail open - allow request if Redis error
		return true, 0
	}

	// If key doesn't exist or count is below limit
	if err == redis.Nil || count < rl.limit {
		// Increment count (equivalent to Lock + increment)
		newCount, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			return true, 0 // Fail open
		}

		// Set expiration for new keys (equivalent to resetCount goroutine)
		if newCount == 1 {
			rl.client.Expire(ctx, key, rl.window)
		}

		return true, 0
	}

	// Get remaining TTL for denied requests
	ttl, err := rl.client.TTL(ctx, key).Result()
	if err != nil {
		ttl = rl.window
	}

	return false, ttl
}
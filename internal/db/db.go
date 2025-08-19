// Package db provides database connection management utilities for PostgreSQL.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

// Default configuration constants
const (
	DefaultConnectTimeout = 30 * time.Second
	DefaultPingTimeout    = 5 * time.Second
	DefaultMaxLifetime    = time.Hour
)

// Config holds database connection configuration
type Config struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	if c.DSN == "" {
		return errors.New("DSN cannot be empty")
	}
	if c.MaxOpenConns < 0 {
		return errors.New("MaxOpenConns cannot be negative")
	}
	if c.MaxIdleConns < 0 {
		return errors.New("MaxIdleConns cannot be negative")
	}
	if c.MaxIdleConns > c.MaxOpenConns && c.MaxOpenConns > 0 {
		return fmt.Errorf("MaxIdleConns (%d) cannot exceed MaxOpenConns (%d)",
			c.MaxIdleConns, c.MaxOpenConns)
	}
	return nil
}

// New creates a new PostgreSQL database connection with the given configuration.
// It configures connection pooling, validates connectivity, and returns a ready-to-use *sql.DB.
//
// The function will:
//   - Open a connection to PostgreSQL
//   - Configure connection pool limits and timeouts
//   - Validate connectivity with a ping
//   - Return an error if any step fails
//
// Example:
//
//	db, err := db.New("postgres://user:pass@localhost/db", 25, 25, "15m")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
func New(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	config := &Config{
		DSN:          dsn,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
		MaxIdleTime:  maxIdleTime,
	}

	return NewWithConfig(config)
}

// NewWithConfig creates a new database connection using the provided Config.
// This is the preferred method for new code as it allows for easier testing and configuration management.
func NewWithConfig(config *Config) (*sql.DB, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Parse idle time duration
	idleDuration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("invalid maxIdleTime format: %w", err)
	}

	// Open database connection
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(idleDuration)
	db.SetConnMaxLifetime(DefaultMaxLifetime)

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), DefaultPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close() // Clean up on failure
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// Ping tests database connectivity with a timeout.
// This is a convenience function for health checks.
func Ping(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultPingTimeout)
	defer cancel()

	return db.PingContext(ctx)
}

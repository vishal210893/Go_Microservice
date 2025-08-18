// Package env provides utilities for reading and parsing environment variables
// with fallback values and type conversion capabilities. It supports common
// configuration patterns for cloud-native applications.
package env

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetString retrieves a string environment variable with fallback support.
// Returns the environment variable value if it exists and is non-empty,
// otherwise returns the provided fallback value.
func GetString(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

// GetPort retrieves a port configuration from environment variables.
// Validates the port format and ensures it starts with ':' for server binding.
// Returns the port string suitable for http.Server.Addr configuration.
func GetPort(key, fallback string) string {
	val := GetString(key, fallback)

	// Ensure port starts with ':' for server binding
	if val != "" && !strings.HasPrefix(val, ":") {
		val = ":" + val
	}

	return val
}

// GetInt retrieves an integer environment variable with fallback support.
// Attempts to parse the environment variable as an integer. If parsing fails
// or the variable doesn't exist, returns the fallback value.
func GetInt(key string, fallback int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		slog.Warn("Failed to parse integer environment variable",
			"key", key,
			"value", val,
			"error", err,
			"using_fallback", fallback)
		return fallback
	}

	return parsed
}

// GetPortInt retrieves a port number as an integer from environment variables.
// Validates the port is within the valid range (1-65535) and returns the
// fallback value if validation fails.
//
// Deprecated: Use GetInt with application-level port validation for better
// error handling and logging.
func GetPortInt(key string, fallback int) int {
	port := GetInt(key, fallback)

	// Validate port range
	if port < 1 || port > 65535 {
		slog.Warn("Invalid port number, using fallback",
			"key", key,
			"invalid_port", port,
			"fallback", fallback)
		return fallback
	}

	return port
}

// GetBool retrieves a boolean environment variable with fallback support.
// Accepts common boolean representations: "true", "false", "1", "0", "yes", "no".
// Case-insensitive parsing with fallback on parsing errors.
func GetBool(key string, fallback bool) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if val == "" {
		return fallback
	}

	switch val {
	case "true", "1", "yes", "on", "enable", "enabled":
		return true
	case "false", "0", "no", "off", "disable", "disabled":
		return false
	default:
		slog.Warn("Failed to parse boolean environment variable",
			"key", key,
			"value", val,
			"using_fallback", fallback)
		return fallback
	}
}

// GetDuration retrieves a duration environment variable with fallback support.
// Accepts Go duration format strings (e.g., "30s", "5m", "1h30m").
// Returns fallback value on parsing errors.
func GetDuration(key string, fallback time.Duration) time.Duration {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(val)
	if err != nil {
		slog.Warn("Failed to parse duration environment variable",
			"key", key,
			"value", val,
			"error", err,
			"using_fallback", fallback)
		return fallback
	}

	return parsed
}

// Required retrieves a required environment variable and panics if it doesn't exist.
// Use this for critical configuration that must be provided at runtime.
// Consider using GetString with validation in application code for better error handling.
func Required(key string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		slog.Error("Required environment variable not set", "key", key)
		panic("missing required environment variable: " + key)
	}
	return val
}

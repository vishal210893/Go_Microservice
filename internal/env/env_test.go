package env

import (
    "testing"
    "time"
)

func TestGetString(t *testing.T) {
    t.Setenv("APP_NAME", "service")
    if got := GetString("APP_NAME", "fallback"); got != "service" {
        t.Fatalf("GetString returned %q, want %q", got, "service")
    }
    // Empty variable should use fallback
    t.Setenv("EMPTY_VAR", "")
    if got := GetString("EMPTY_VAR", "fallback"); got != "fallback" {
        t.Fatalf("GetString fallback = %q, want %q", got, "fallback")
    }
}

func TestGetPort(t *testing.T) {
    t.Setenv("PORT_ONLY", "8081")
    if got := GetPort("PORT_ONLY", ":8080"); got != ":8081" {
        t.Fatalf("GetPort(8081) = %q, want %q", got, ":8081")
    }
    t.Setenv("PORT_WITH_COLON", ":9090")
    if got := GetPort("PORT_WITH_COLON", ":8080"); got != ":9090" {
        t.Fatalf("GetPort(:9090) = %q, want %q", got, ":9090")
    }
}

func TestGetInt(t *testing.T) {
    t.Setenv("WORKERS", "42")
    if got := GetInt("WORKERS", 1); got != 42 {
        t.Fatalf("GetInt valid = %d, want %d", got, 42)
    }
    t.Setenv("WORKERS", "abc")
    if got := GetInt("WORKERS", 7); got != 7 {
        t.Fatalf("GetInt fallback on invalid = %d, want %d", got, 7)
    }
}

func TestGetBool(t *testing.T) {
    t.Setenv("FEATURE_X", "true")
    if !GetBool("FEATURE_X", false) {
        t.Fatalf("GetBool(true) = false, want true")
    }
    t.Setenv("FEATURE_X", "0")
    if GetBool("FEATURE_X", true) {
        t.Fatalf("GetBool(0) = true, want false")
    }
    t.Setenv("FEATURE_X", "maybe")
    if !GetBool("FEATURE_X", true) { // falls back
        t.Fatalf("GetBool(invalid) should fallback to true")
    }
}

func TestGetDuration(t *testing.T) {
    t.Setenv("TIMEOUT", "15m")
    if got := GetDuration("TIMEOUT", time.Minute); got != 15*time.Minute {
        t.Fatalf("GetDuration valid = %v, want %v", got, 15*time.Minute)
    }
    t.Setenv("TIMEOUT", "notaduration")
    if got := GetDuration("TIMEOUT", 30*time.Second); got != 30*time.Second {
        t.Fatalf("GetDuration fallback = %v, want %v", got, 30*time.Second)
    }
}

func TestRequiredPanicsWhenMissing(t *testing.T) {
    t.Setenv("REQUIRED_X", "")
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("Required did not panic for missing var")
        }
    }()
    _ = Required("REQUIRED_X")
}


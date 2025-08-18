package env

import (
	"os"
	"strconv"
)

func GetPort(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetPortInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	atoi, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return atoi
}

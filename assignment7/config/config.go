package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort          string
	JWTSecret         string
	DBDSN             string
	RateLimitRequests int
	RateLimitWindowS  int
}

func New() *Config {
	rateLimitRequests := readInt("RATE_LIMIT_REQUESTS", 20)
	rateLimitWindow := readInt("RATE_LIMIT_WINDOW_SECONDS", 60)

	return &Config{
		HTTPPort:          readString("HTTP_PORT", "8080"),
		JWTSecret:         readString("JWT_SECRET", "change-me-in-production"),
		DBDSN:             readString("DB_DSN", "practice7.db"),
		RateLimitRequests: rateLimitRequests,
		RateLimitWindowS:  rateLimitWindow,
	}
}

func readString(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func readInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}

	return n
}

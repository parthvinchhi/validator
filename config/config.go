// Package config reads the (very small) configuration from environment variables.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port               string
	Environment        string
	APIKey             string // empty = API key check disabled
	CORSAllowedOrigins []string
	TrustedProxies     []string
	RateLimitPerMinute int   // 0 = disabled
	MaxBodyBytes       int64 // max size of a JSON request body
}

func Load() Config {
	cfg := Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		APIKey:             os.Getenv("API_KEY"),
		CORSAllowedOrigins: splitList(os.Getenv("CORS_ALLOWED_ORIGINS")),
		TrustedProxies:     splitList(os.Getenv("TRUSTED_PROXIES")),
		RateLimitPerMinute: 300,
		MaxBodyBytes:       4096,
	}

	if raw := os.Getenv("RATE_LIMIT_PER_MINUTE"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			log.Println("Invalid RATE_LIMIT_PER_MINUTE, using default of 300")
		} else {
			cfg.RateLimitPerMinute = n
		}
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// splitList turns "a, b ,c" into ["a" "b" "c"]. Empty input gives nil.
func splitList(raw string) []string {
	var items []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}

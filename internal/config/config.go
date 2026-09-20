// Package config loads settings from environment variables.
package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      string
	APIKeys   []string // empty = auth disabled (dev only)
	RateRPS   float64  // requests/second per API key (or IP when auth is off)
	RateBurst int
	MaxBatch  int
}

func Load() Config {
	var keys []string
	for _, k := range strings.Split(os.Getenv("API_KEYS"), ",") {
		if k = strings.TrimSpace(k); k != "" {
			keys = append(keys, k)
		}
	}
	return Config{
		Port:      getenv("PORT", "8080"),
		APIKeys:   keys,
		RateRPS:   getfloat("RATE_LIMIT_RPS", 20),
		RateBurst: getint("RATE_LIMIT_BURST", 40),
		MaxBatch:  getint("MAX_BATCH", 100),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getint(k string, def int) int {
	if n, err := strconv.Atoi(os.Getenv(k)); err == nil && n > 0 {
		return n
	}
	return def
}

func getfloat(k string, def float64) float64 {
	if f, err := strconv.ParseFloat(os.Getenv(k), 64); err == nil && f > 0 {
		return f
	}
	return def
}

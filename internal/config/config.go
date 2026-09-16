package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	DatabaseDirectURL    string
	OpenRouterAPIKey     string
	OpenRouterBaseURL    string
	OpenRouterModel      string
	OpenRouterFallback   string
	OpenRouterReferer    string
	OpenRouterAppName    string
	ZeroCostMode         bool
	ModelTimeout         time.Duration
	ModelRetries         int
	ModelMaxOutputTokens int
	CORSOrigins          string
}

func Load() Config {
	loadDotEnv(".env")
	return Config{
		Port: env("PORT", "8080"), DatabaseURL: os.Getenv("DATABASE_URL"), DatabaseDirectURL: os.Getenv("DATABASE_URL_DIRECT"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"), OpenRouterBaseURL: env("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		OpenRouterModel: env("OPENROUTER_MODEL", "openrouter/free"), OpenRouterFallback: os.Getenv("OPENROUTER_FALLBACK_MODEL"),
		OpenRouterReferer: env("OPENROUTER_HTTP_REFERER", "http://localhost:3000"), OpenRouterAppName: env("OPENROUTER_APP_NAME", "MuleLab"),
		ZeroCostMode: boolEnv("ZERO_COST_MODE", true), ModelTimeout: time.Duration(intEnv("MODEL_REQUEST_TIMEOUT_SECONDS", 25)) * time.Second,
		ModelRetries: intEnv("MODEL_MAX_RETRIES", 1), ModelMaxOutputTokens: intEnv("MODEL_MAX_OUTPUT_TOKENS", 1200), CORSOrigins: env("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if ok && os.Getenv(strings.TrimSpace(key)) == "" {
			_ = os.Setenv(strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"`))
		}
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func intEnv(k string, d int) int {
	if v, err := strconv.Atoi(os.Getenv(k)); err == nil && v > 0 {
		return v
	}
	return d
}
func boolEnv(k string, d bool) bool {
	if v, err := strconv.ParseBool(os.Getenv(k)); err == nil {
		return v
	}
	return d
}

package config

import (
	"strings"
	"time"

	"github.com/SeakMengs/yato-cdn/internal/env"
)

type Config struct {
	Port        string
	ENV         string
	DB          DatabaseConfig
	RateLimiter RateLimiterConfig
	CDN         CDN
}

type CDN struct {
	IsCDN  bool
	Region string
}

type RateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

type DatabaseConfig struct {
	DB_HOST      string
	DB_PORT      string
	DB_DATABASE  string
	DB_USERNAME  string
	DB_PASSWORD  string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.ENV, "production")
}

func GetConfig() Config {
	rateLimiteTimeFrame, err := time.ParseDuration(env.GetString("RATE_LIMIT_TIME_FRAME", "1m"))
	if err != nil {
		rateLimiteTimeFrame = 60 * time.Second
	}

	return Config{
		Port: env.GetString("PORT", "8080"),
		ENV:  env.GetString("ENV", "development"),
		DB: DatabaseConfig{
			DB_HOST:      env.GetString("DB_HOST", "127.0.0.1"),
			DB_PORT:      env.GetString("DB_PORT", "5432"),
			DB_USERNAME:  env.GetString("DB_USERNAME", "root"),
			DB_PASSWORD:  env.GetString("DB_PASSWORD", ""),
			DB_DATABASE:  env.GetString("DB_DATABASE", "database_name"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		// By default if not specified, we allow 5000 requests per minute on all routes
		RateLimiter: RateLimiterConfig{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMIT_REQUESTS_PER_TIME_FRAME", 5000),
			TimeFrame:            rateLimiteTimeFrame,
			Enabled:              env.GetBool("RATE_LIMIT_ENABLED", true),
		},
		CDN: CDN{
			IsCDN:  env.GetBool("IS_CDN", false),
			Region: env.GetString("REGION", "universal"),
		},
	}
}

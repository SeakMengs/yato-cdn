package ratelimiter

import (
	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/util"
	"go.uber.org/zap"
)

func NewRateLimiter(cfg config.RateLimiterConfig, logger *zap.SugaredLogger) *FixedWindowRateLimiter {
	// For unit test
	if logger == nil {
		logger = util.NewLogger()
	}

	return NewFixedWindowLimiter(cfg, logger)
}

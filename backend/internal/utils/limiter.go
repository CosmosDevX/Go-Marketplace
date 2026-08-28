package utils

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

func ActivateRateLimiter(ctx context.Context, w http.ResponseWriter, r *http.Request, key string, rateLimiter *redis_rate.Limiter, limit redis_rate.Limit) error {
	ip := GetIP(r)

	res, err := rateLimiter.Allow(ctx, key+ip, limit)
	if err != nil {
		WriteError(w, fmt.Errorf("rate limit allow: %w", domain.ErrInternalServerError))
		return fmt.Errorf("rate limit: %w", domain.ErrInternalServerError)
	}

	if res.Allowed == 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", res.RetryAfter.Seconds()))
		WriteError(w, fmt.Errorf("not allowed: %w", domain.ErrTooManyRequests))
		return fmt.Errorf("rate limit: %w", domain.ErrTooManyRequests)
	}

	return nil
}

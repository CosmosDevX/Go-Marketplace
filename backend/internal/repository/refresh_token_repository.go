// Package repository
package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository struct {
	redisClient *redis.Client
}

func NewRefreshTokenRepository(redisClient *redis.Client) RefreshTokenRepository {
	return RefreshTokenRepository{
		redisClient: redisClient,
	}
}

func (r RefreshTokenRepository) Set(ctx context.Context, userID, refreshToken, prefix string, ttl time.Duration) error {
	err := r.redisClient.Set(ctx, userID+prefix, refreshToken, ttl).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("set refresh token by userID %s: %w", userID, domain.ErrTimeout)
		}

		return fmt.Errorf("set refresh token by userID %s: %w", userID, err)
	}

	return nil
}

func (r RefreshTokenRepository) Delete(ctx context.Context, userID, prefix string) error {
	err := r.redisClient.Del(ctx, userID+prefix).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete refresh token by userID %s: %w", userID, domain.ErrTimeout)
		}

		return fmt.Errorf("delete refresh token by userID %s: %w", userID, err)
	}

	return nil
}

func (r RefreshTokenRepository) Get(ctx context.Context, userID, prefix string) (string, error) {
	cmd := r.redisClient.Get(ctx, userID+prefix)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("get refresh token by userID %s: %w", userID, domain.ErrTimeout)
		}
		if errors.Is(cmd.Err(), redis.Nil) {
			return "", fmt.Errorf("get refresh token by userID %s: %w", userID, domain.ErrNotFound)
		}

		return "", fmt.Errorf("get refresh token by userID %s: %w", userID, cmd.Err())
	}

	return cmd.Val(), nil
}

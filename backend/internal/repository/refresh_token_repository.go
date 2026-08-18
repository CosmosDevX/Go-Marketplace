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

func (r RefreshTokenRepository) Set(ctx context.Context, refreshToken, userID, prefix string, ttl time.Duration) error {
	err := r.redisClient.Set(ctx, refreshToken+prefix, userID, ttl).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("set userID by refresh token: %w", domain.ErrTimeout)
		}

		return fmt.Errorf("set userID by refresh token: %w", err)
	}

	return nil
}

func (r RefreshTokenRepository) Delete(ctx context.Context, refreshToken, prefix string) error {
	err := r.redisClient.Del(ctx, refreshToken+prefix).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete userID by refresh token: %w", domain.ErrTimeout)
		}

		return fmt.Errorf("delete userID by refresh token: %w", err)
	}

	return nil
}

func (r RefreshTokenRepository) Get(ctx context.Context, refreshToken, prefix string) (string, error) {
	cmd := r.redisClient.Get(ctx, refreshToken+prefix)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("get userID by refresh token: %w", domain.ErrTimeout)
		}
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("get userID by refresh token: %w", domain.ErrNotFound)
		}

		return "", fmt.Errorf("get userID by refresh token: %w", err)
	}

	return cmd.Val(), nil
}

func (r RefreshTokenRepository) GetAndDelete(ctx context.Context, refreshToken, prefix string) (string, error) {
	cmd := r.redisClient.GetDel(ctx, refreshToken+prefix)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("delete userID by refresh token: %w", domain.ErrTimeout)
		}
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("delete userID by refresh token: %w", domain.ErrNotFound)
		}

		return "", fmt.Errorf("delete userID by refresh token: %w", err)
	}

	return cmd.Val(), nil
}

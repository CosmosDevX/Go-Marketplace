// Package repository
package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
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

func (r RefreshTokenRepository) Set(ctx context.Context, userID, refreshToken, prefix string, ttl time.Duration) *domain.DomainError {
	err := r.redisClient.Set(ctx, userID+prefix, refreshToken, ttl).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return &domain.DomainError{Code: constants.SaveError, Message: "error during save refresh token"}
	}

	return nil
}

func (r RefreshTokenRepository) Delete(ctx context.Context, userID, prefix string) *domain.DomainError {
	err := r.redisClient.Del(ctx, userID+prefix).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return &domain.DomainError{Code: constants.DeleteError, Message: "error during delete refresh token"}
	}

	return nil
}

func (r RefreshTokenRepository) Get(ctx context.Context, userID, prefix string) (string, *domain.DomainError) {
	cmd := r.redisClient.Get(ctx, userID+prefix)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(cmd.Err(), redis.Nil) {
			return "", nil
		}
		return "", &domain.DomainError{Code: constants.FindError, Message: "no matches found for this username"}
	}

	return cmd.Val(), nil
}

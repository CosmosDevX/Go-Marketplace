// Package repository
package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"

	"github.com/redis/go-redis/v9"
)

const accessTokenBlacklistKeyPrefix = ":accessTokensBlacklist"

type AccessTokenRepository struct {
	redisClient *redis.Client
}

func NewAccessTokenRepository(redisClient *redis.Client) AccessTokenRepository {
	return AccessTokenRepository{
		redisClient: redisClient,
	}
}

func (r AccessTokenRepository) Set(ctx context.Context, accessTokenHash string) error {
	err := r.redisClient.Set(ctx, accessTokenHash+accessTokenBlacklistKeyPrefix, "", config.AccessTokenTTL).Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("set access token hash in blacklist: %w", domain.ErrTimeout)
		}

		return fmt.Errorf("set access token hash in blacklist: %w", err)
	}

	return nil
}

func (r AccessTokenRepository) Exists(ctx context.Context, accessTokenHash string) (bool, error) {
	cmd := r.redisClient.Get(ctx, accessTokenHash+accessTokenBlacklistKeyPrefix)
	err := cmd.Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return false, fmt.Errorf("check access token hash in blacklist: %w", domain.ErrTimeout)
		}
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("check access token hash in blacklist: %w", err)
	}

	return true, nil
}

// Package middleware
package middleware

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/logger"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"
	"strings"
)

type JWTServiceInterface interface {
	ParseAccessToken(tokenString string) (authorization.AccessTokenClaims, error)
}

type AccessTokenBlacklistChecker interface {
	Exists(ctx context.Context, accessTokenHash string) (bool, error)
}

type AuthMiddleware struct {
	jwtService            JWTServiceInterface
	tokenBlacklistChecker AccessTokenBlacklistChecker
}

func NewAuthMiddleware(jwtService JWTServiceInterface, tokenBlacklistChecker AccessTokenBlacklistChecker) AuthMiddleware {
	return AuthMiddleware{
		jwtService:            jwtService,
		tokenBlacklistChecker: tokenBlacklistChecker,
	}
}

type UserContextKey struct{}
type UserContext struct {
	UserID      int
	Roles       []string
	AccessToken string
}

func UserFromContext(ctx context.Context) (UserContext, error) {
	ctxValue := ctx.Value(UserContextKey{})
	userContext, ok := ctxValue.(UserContext)
	if !ok {
		return UserContext{}, fmt.Errorf("get user from context: %w", domain.ErrParse)
	}

	return userContext, nil
}

func (m AuthMiddleware) ProtectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(ctx, w, fmt.Errorf("auth header is empty: %w", domain.ErrUnauthorized))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(ctx, w, fmt.Errorf("invalid auth scheme type: %w", domain.ErrUnauthorized))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			utils.WriteError(ctx, w, fmt.Errorf("token is empty: %w", domain.ErrUnauthorized))
			return
		}

		claims, err := m.jwtService.ParseAccessToken(tokenString)
		if err != nil {
			utils.WriteError(ctx, w, fmt.Errorf("auth token: %w", domain.ErrUnauthorized))
			return
		}

		isExists, err := m.tokenBlacklistChecker.Exists(ctx, utils.HashToken(tokenString))
		if err != nil {
			utils.WriteError(ctx, w, err)
			return
		}
		if isExists {
			utils.WriteError(ctx, w, fmt.Errorf("access token in blacklist: %w", domain.ErrUnauthorized))
			return
		}

		ctx = context.WithValue(ctx, UserContextKey{}, UserContext{
			UserID:      claims.UserID,
			Roles:       claims.Roles,
			AccessToken: tokenString,
		})
		ctx = logger.WithUserID(ctx, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

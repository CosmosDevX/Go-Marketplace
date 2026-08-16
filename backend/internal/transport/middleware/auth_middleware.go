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

type UsernameContextKey struct{}
type UserIDContextKey struct{}

type AuthMiddleware struct {
	jwtService authorization.JWTServiceInterface
}

func NewAuthMiddleware(jwtService authorization.JWTServiceInterface) AuthMiddleware {
	return AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m AuthMiddleware) ProtectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, fmt.Errorf("auth header is empty: %w", domain.ErrUnauthorized))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(w, fmt.Errorf("invalid auth scheme type: %w", domain.ErrUnauthorized))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			utils.WriteError(w, fmt.Errorf("token is empty: %w", domain.ErrUnauthorized))
			return
		}

		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			utils.WriteError(w, fmt.Errorf("auth token: %w", domain.ErrUnauthorized))
			return
		}

		ctx = context.WithValue(ctx, UsernameContextKey{}, claims.Username)
		ctx = context.WithValue(ctx, UserIDContextKey{}, claims.UserID)
		ctx = logger.WithUserID(ctx, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

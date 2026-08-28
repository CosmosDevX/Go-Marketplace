// Package handler
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type AuthService interface {
	Auth(ctx context.Context, input authorization.LoginInput) (authorization.AuthResult, error)
	Refresh(ctx context.Context, oldRefreshToken string) (authorization.AuthResult, error)
	Logout(ctx context.Context, refreshToken, accessToken string) error
}

const (
	refreshTokenKey     = "refresh_token"
	refreshTokenMaxAge  = 3600 * 24 * 7
	authRateLimitKey    = "auth"
	refreshRateLimitKey = "refresh"
)

type AuthHandler struct {
	authService AuthService
	rateLimiter redis_rate.Limiter
}

func NewAuthHandler(authService AuthService, rateLimiter redis_rate.Limiter) AuthHandler {
	return AuthHandler{
		authService: authService,
		rateLimiter: rateLimiter,
	}
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, authRateLimitKey, &h.rateLimiter, redis_rate.PerMinute(5)); err != nil {
		return
	}

	var dto dto.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("login dto deserialize: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("login dto validate: %w", domain.ErrValidation))
		return
	}

	authResult, err := h.authService.Auth(ctx, authorization.LoginInput{
		Username: dto.Username,
		Password: dto.Password,
	})
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]any{"access_token": authResult.AccessToken, "roles": authResult.Roles})
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, refreshRateLimitKey, &h.rateLimiter, redis_rate.PerMinute(3)); err != nil {
		return
	}

	tokenCookie, err := r.Cookie(refreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteError(ctx, w, fmt.Errorf("refresh token not exists: %w", domain.ErrUnauthorized))
		return
	}

	authResult, err := h.authService.Refresh(ctx, tokenCookie.Value)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))
	utils.WriteJSON(w, map[string]any{"access_token": authResult.AccessToken, "roles": authResult.Roles})
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(refreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteError(ctx, w, fmt.Errorf("refresh token not exists: %w", domain.ErrUnauthorized))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}
	if err := h.authService.Logout(ctx, tokenCookie.Value, user.AccessToken); err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    "",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	utils.WriteJSON(w, map[string]string{"message": "logout successful"})
}

func (h AuthHandler) newRefreshTokenCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     refreshTokenKey,
		Value:    refreshToken,
		MaxAge:   refreshTokenMaxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

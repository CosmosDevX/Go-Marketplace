// Package handler
package handler

import (
	"context"
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
	Auth(ctx context.Context, userDTO dto.LoginDTO) (*authorization.AuthResult, error)
	Refresh(ctx context.Context, oldRefreshToken string) (*authorization.AuthResult, error)
	Logout(ctx context.Context, userID int) error
}

const (
	refreshTokenKey    = "refresh_token"
	refreshTokenMaxAge = 3600 * 24 * 7
	authRateLimitKey   = "auth"
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

func (h AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ActivateRateLimiter(ctx, w, r, authRateLimitKey, &h.rateLimiter, redis_rate.PerMinute(5)); err != nil {
		return
	}

	var loginDTO dto.LoginDTO
	if err := utils.Deserialize(r.Body, &loginDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("login dto deserialize: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(loginDTO); err != nil {
		utils.WriteError(w, fmt.Errorf("login dto validate: %w", domain.ErrValidation))
		return
	}

	authResult, err := h.authService.Auth(ctx, loginDTO)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(refreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteError(w, fmt.Errorf("refresh token not exists: %w", domain.ErrUnauthorized))
		return
	}

	authResult, err := h.authService.Refresh(ctx, tokenCookie.Value)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(refreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteJSON(w, map[string]string{"message": "refresh token not exists"})
		return
	}

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, fmt.Errorf("user id parse: %w", domain.ErrParse))
		return
	}

	if err := h.authService.Logout(ctx, userID); err != nil {
		utils.WriteError(w, err)
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

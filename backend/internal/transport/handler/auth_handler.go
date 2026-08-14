// Package handler
package handler

import (
	"context"
	"myapp/internal/constants"
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
	Auth(ctx context.Context, userDTO dto.LoginDTO) (*authorization.AuthResult, *domain.DomainError)
	Refresh(ctx context.Context, oldRefreshToken string) (*authorization.AuthResult, *domain.DomainError)
	Logout(ctx context.Context, userID int) *domain.DomainError
}

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

	if err := utils.ActivateRateLimiter(ctx, w, r, constants.AuthRateLimitKey, &h.rateLimiter, redis_rate.PerMinute(5)); err != nil {
		return
	}

	var loginDTO dto.LoginDTO
	if err := utils.Deserialize(r.Body, &loginDTO); err != nil {
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing user"))
		return
	}

	if err := validator.Struct(loginDTO); err != nil {
		utils.WriteError(w, *domain.NewValidationError(err.Error()))
		return
	}

	authResult, domainErr := h.authService.Auth(ctx, loginDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteError(w, *domain.NewDomainError(constants.InvalidTokenError, "refresh token not exists"))
		return
	}

	authResult, domainErr := h.authService.Refresh(ctx, tokenCookie.Value)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteJSON(w, map[string]string{"message": "refresh token not exists"})
		return
	}

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewDomainError(constants.ParseError, "error during parse user id"))
		return
	}

	if domainErr := h.authService.Logout(ctx, userID); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenKey,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.WriteJSON(w, map[string]string{"message": "logout successful"})
}

func (h AuthHandler) newRefreshTokenCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     constants.RefreshTokenKey,
		Value:    refreshToken,
		MaxAge:   constants.RefreshTokenMaxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

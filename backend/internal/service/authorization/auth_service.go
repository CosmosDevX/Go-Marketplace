package authorization

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string
	Password string
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type JWTServiceInterface interface {
	GenerateAccessToken(userID int, roles []string, expiresAt time.Duration) (string, error)
	GenerateRefreshToken() (string, error)
	ParseAccessToken(tokenString string) (AccessTokenClaims, error)
}

type UserRolesGetter interface {
	ListRolesByUserID(ctx context.Context, userID int) ([]string, error)
}

type UserGetter interface {
	GetByName(ctx context.Context, username string) (domain.User, error)
}

type RefreshTokenRepository interface {
	Set(ctx context.Context, refreshToken, userID, prefix string, ttl time.Duration) error
	Delete(ctx context.Context, refreshToken, prefix string) error
	Get(ctx context.Context, refreshToken, prefix string) (string, error)
	GetAndDelete(ctx context.Context, refreshToken, prefix string) (string, error)
}

const (
	accessTokenExpiresAt  = 15 * time.Minute
	refreshTokenExpiresAt = 24 * time.Hour * 7
	tokenWhiteListPrefix  = ":tokensWhiteList"
)

type AuthService struct {
	userGetter             UserGetter
	userRolesGetter        UserRolesGetter
	refreshTokenRepository RefreshTokenRepository
	jwtService             JWTServiceInterface
}

func NewAuthService(userGetter UserGetter, userRolesGetter UserRolesGetter, refreshTokenRepo RefreshTokenRepository, jwtService JWTServiceInterface) AuthService {
	return AuthService{
		userGetter:             userGetter,
		userRolesGetter:        userRolesGetter,
		refreshTokenRepository: refreshTokenRepo,
		jwtService:             jwtService,
	}
}

func (s AuthService) Auth(ctx context.Context, input LoginInput) (AuthResult, error) {
	user, err := s.userGetter.GetByName(ctx, input.Username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, fmt.Errorf("auth: %w", domain.ErrUnauthorized)
		}
		return AuthResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, fmt.Errorf("password compare: %w", domain.ErrUnauthorized)
	}

	authResult, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	if err := s.refreshTokenRepository.Set(ctx, authResult.RefreshToken, strconv.Itoa(user.ID), tokenWhiteListPrefix, refreshTokenExpiresAt); err != nil {
		return AuthResult{}, err
	}

	return authResult, nil
}

func (s AuthService) Refresh(ctx context.Context, oldRefreshToken string) (AuthResult, error) {
	value, err := s.refreshTokenRepository.GetAndDelete(ctx, oldRefreshToken, tokenWhiteListPrefix)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, fmt.Errorf("refresh token get: %w", domain.ErrUnauthorized)
		}
		return AuthResult{}, err
	}

	userID, err := strconv.Atoi(value)
	if err != nil {
		return AuthResult{}, fmt.Errorf("userID parse: %w", domain.ErrParse)
	}

	authResult, err := s.generateTokenPair(ctx, userID)
	if err != nil {
		return AuthResult{}, err
	}

	if err := s.refreshTokenRepository.Set(ctx, authResult.RefreshToken, strconv.Itoa(userID), tokenWhiteListPrefix, refreshTokenExpiresAt); err != nil {
		return AuthResult{}, err
	}

	return authResult, nil
}

func (s AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepository.Delete(ctx, refreshToken, tokenWhiteListPrefix); err != nil {
		return err
	}

	return nil
}

func (s AuthService) generateTokenPair(ctx context.Context, userID int) (AuthResult, error) {
	roles, err := s.userRolesGetter.ListRolesByUserID(ctx, userID)
	if err != nil {
		return AuthResult{}, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(userID, roles, accessTokenExpiresAt)
	if err != nil {
		return AuthResult{}, fmt.Errorf("access token generate: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return AuthResult{}, fmt.Errorf("refresh token generate: %w", err)
	}

	return AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

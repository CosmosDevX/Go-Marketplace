package authorization

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type UserGetter interface {
	GetByName(ctx context.Context, username string) (*domain.User, error)
}

type RefreshTokenRepository interface {
	Set(ctx context.Context, userID, refreshToken, prefix string, ttl time.Duration) error
	Delete(ctx context.Context, userID, prefix string) error
	Get(ctx context.Context, userID, prefix string) (string, error)
}

const (
	accessTokenExpiresAt  = 15 * time.Minute
	refreshTokenExpiresAt = 24 * time.Hour * 7
	tokenWhiteListPrefix  = ":tokensWhiteList"
)

type AuthService struct {
	userGetter             UserGetter
	refreshTokenRepository RefreshTokenRepository
	jwtService             JWTServiceInterface
}

func NewAuthService(userGetter UserGetter, refreshTokenRepo RefreshTokenRepository, jwtService JWTServiceInterface) AuthService {
	return AuthService{
		userGetter:             userGetter,
		refreshTokenRepository: refreshTokenRepo,
		jwtService:             jwtService,
	}
}

func (s AuthService) Auth(ctx context.Context, userDTO dto.LoginDTO) (*AuthResult, error) {
	user, err := s.userGetter.GetByName(ctx, userDTO.Username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("auth: %w", domain.ErrUnauthorized)
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDTO.Password)); err != nil {
		return nil, fmt.Errorf("password compare: %w", domain.ErrUnauthorized)
	}

	authResult, err := s.generateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Set(ctx, strconv.Itoa(user.ID), authResult.RefreshToken, tokenWhiteListPrefix, refreshTokenExpiresAt); err != nil {
		return nil, err
	}

	return authResult, nil
}

func (s AuthService) Refresh(ctx context.Context, oldRefreshToken string) (*AuthResult, error) {
	claims, err := s.jwtService.ParseToken(oldRefreshToken)
	if err != nil {
		return nil, err
	}

	dbRefreshToken, err := s.refreshTokenRepository.Get(ctx, strconv.Itoa(claims.UserID), tokenWhiteListPrefix)
	if err != nil {
		return nil, err
	}
	if oldRefreshToken != dbRefreshToken {
		return nil, fmt.Errorf("refresh tokens comparing: %w", domain.ErrUnauthorized)
	}

	authResult, err := s.generateTokenPair(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Set(ctx, strconv.Itoa(claims.UserID), authResult.RefreshToken, tokenWhiteListPrefix, refreshTokenExpiresAt); err != nil {
		return nil, err
	}

	return authResult, nil
}

func (s AuthService) Logout(ctx context.Context, userID int) error {
	if _, err := s.refreshTokenRepository.Get(ctx, strconv.Itoa(userID), tokenWhiteListPrefix); err != nil {
		return err
	}

	if err := s.refreshTokenRepository.Delete(ctx, strconv.Itoa(userID), tokenWhiteListPrefix); err != nil {
		return err
	}

	return nil
}

func (s AuthService) generateTokenPair(userID int, username string) (*AuthResult, error) {
	accessToken, jwtError := s.jwtService.GenerateToken(userID, username, accessTokenExpiresAt)
	if jwtError != nil {
		return nil, fmt.Errorf("access token generate: %w", jwtError)
	}

	refreshToken, jwtError := s.jwtService.GenerateToken(userID, username, refreshTokenExpiresAt)
	if jwtError != nil {
		return nil, fmt.Errorf("refresh token generate: %w", jwtError)
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

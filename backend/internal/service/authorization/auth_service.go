package authorization

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/logger"
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
	GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError)
}

type RefreshTokenRepository interface {
	Set(ctx context.Context, userID, refreshToken, prefix string, ttl time.Duration) *domain.DomainError
	Delete(ctx context.Context, userID, prefix string) *domain.DomainError
	Get(ctx context.Context, userID, prefix string) (string, *domain.DomainError)
}

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

func (s AuthService) Auth(ctx context.Context, userDTO dto.LoginDTO) (*AuthResult, *domain.DomainError) {
	log := logger.FromContext(ctx)

	user, err := s.userGetter.GetUserByName(ctx, userDTO.Username)
	if err != nil {
		log.Warn("auth failed: user not found", "username", userDTO.Username, "code", err.Code)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		log.Warn("auth failed: invalid password", "username", userDTO.Username, "user_id", user.ID)
		return nil, &domain.DomainError{Code: constants.InvalidPassword, Message: "passwords do not match"}
	}

	authResult, domainErr := s.generateTokenPair(user.ID, user.Username)
	if domainErr != nil {
		log.Error("auth failed: token generation", "username", userDTO.Username, "user_id", user.ID, "code", domainErr.Code)
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(ctx, strconv.Itoa(user.ID), authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt); domainErr != nil {
		log.Error("auth failed: store refresh token", "username", userDTO.Username, "user_id", user.ID, "code", domainErr.Code)
		return nil, domainErr
	}

	log.Info("auth success", "username", user.Username, "user_id", user.ID)
	return authResult, nil
}

func (s AuthService) Refresh(ctx context.Context, oldRefreshToken string) (*AuthResult, *domain.DomainError) {
	log := logger.FromContext(ctx)

	claims, domainErr := s.jwtService.ParseToken(oldRefreshToken)
	if domainErr != nil {
		log.Warn("refresh failed: invalid token", "code", domainErr.Code)
		return nil, domainErr
	}

	dbRefreshToken, domainErr := s.refreshTokenRepository.Get(ctx, strconv.Itoa(claims.UserID), constants.TokenWhiteListPrefix)
	if domainErr != nil {
		log.Warn("refresh failed: whitelist lookup", "user_id", claims.UserID, "code", domainErr.Code)
		return nil, domainErr
	}
	if oldRefreshToken != dbRefreshToken {
		log.Warn("refresh failed: token mismatch (possible reuse/theft)", "user_id", claims.UserID)
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	authResult, domainErr := s.generateTokenPair(claims.UserID, claims.Username)
	if domainErr != nil {
		log.Error("refresh failed: token generation", "user_id", claims.UserID, "code", domainErr.Code)
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(ctx, strconv.Itoa(claims.UserID), authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt); domainErr != nil {
		log.Error("refresh failed: store new token", "user_id", claims.UserID, "code", domainErr.Code)
		return nil, domainErr
	}

	log.Info("refresh success", "user_id", claims.UserID, "username", claims.Username)
	return authResult, nil
}

func (s AuthService) Logout(ctx context.Context, userID int) *domain.DomainError {
	log := logger.FromContext(ctx)

	if _, domainErr := s.refreshTokenRepository.Get(ctx, strconv.Itoa(userID), constants.TokenWhiteListPrefix); domainErr != nil {
		log.Warn("logout failed: whitelist lookup", "user_id", userID, "code", domainErr.Code)
		return domainErr
	}

	if domainErr := s.refreshTokenRepository.Delete(ctx, strconv.Itoa(userID), constants.TokenWhiteListPrefix); domainErr != nil {
		log.Error("logout failed: delete whitelist", "user_id", userID, "code", domainErr.Code)
		return domainErr
	}

	log.Info("logout success", "user_id", userID)
	return nil
}

func (s AuthService) generateTokenPair(userID int, username string) (*AuthResult, *domain.DomainError) {
	accessToken, jwtError := s.jwtService.GenerateToken(userID, username, constants.AccessTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.AccessTokenError, Message: "error during the access token generating"}
	}

	refreshToken, jwtError := s.jwtService.GenerateToken(userID, username, constants.RefreshTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.RefreshTokenError, Message: "error during the refresh token generating"}
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

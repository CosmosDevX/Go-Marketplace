// Package authorization
package authorization

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"myapp/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID int
	Roles  []string
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secretKey string) JWTService {
	return JWTService{
		secretKey: []byte(secretKey),
	}
}

func (s JWTService) GenerateAccessToken(userID int, roles []string, expiresAt time.Duration) (string, error) {
	claims := AccessTokenClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myapp",
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("access token sign: %w", err)
	}

	return tokenString, nil
}

func (s JWTService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("refresh token generate: %w", domain.ErrInternalServerError)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s JWTService) ParseAccessToken(tokenString string) (AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect encrypt method")
		}

		return s.secretKey, nil
	})

	if err != nil {
		return AccessTokenClaims{}, fmt.Errorf("token parse: %w", domain.ErrUnauthorized)
	}

	if userClaims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return *userClaims, nil
	}

	return AccessTokenClaims{}, fmt.Errorf("token parse: %w", domain.ErrUnauthorized)
}

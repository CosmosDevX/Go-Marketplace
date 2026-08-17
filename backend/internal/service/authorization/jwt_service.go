// Package authorization
package authorization

import (
	"errors"
	"fmt"
	"time"

	"myapp/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   int
	Username string
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

func (s JWTService) GenerateToken(userID int, username string, expiresAt time.Duration) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		Username: username,
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
		return "", fmt.Errorf("token sign: %w", err)
	}

	return tokenString, nil
}

func (s JWTService) ParseToken(tokenString string) (UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect encrypt method")
		}

		return s.secretKey, nil
	})

	if err != nil {
		return UserClaims{}, fmt.Errorf("token parse: %w", domain.ErrUnauthorized)
	}

	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return *userClaims, nil
	}

	return UserClaims{}, fmt.Errorf("token parse: %w", domain.ErrUnauthorized)
}

// Package domain
package domain

import "fmt"

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Email        string
}

func NewUser(username, passwordHash, email string) (User, error) {
	if username == "" {
		return User{}, fmt.Errorf("username: %w", ErrValidation)
	}
	if email == "" {
		return User{}, fmt.Errorf("email: %w", ErrValidation)
	}
	if passwordHash == "" {
		return User{}, fmt.Errorf("password: %w", ErrValidation)
	}

	return User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
	}, nil
}

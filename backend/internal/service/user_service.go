package service

import (
	"context"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"

	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	Create(ctx context.Context, u domain.User) (int, error)
}

type UserService struct {
	userCreator UserCreator
}

func NewUserService(userCreator UserCreator) UserService {
	return UserService{
		userCreator: userCreator,
	}
}

func (s UserService) Create(ctx context.Context, dto dto.CreateUserDTO) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	if err != nil {
		return 0, fmt.Errorf("password hashing: %w", err)
	}

	passwordHash := string(hashedPassword)

	user, err := domain.NewUser(dto.Username, passwordHash, dto.Email)
	if err != nil {
		return 0, err
	}

	id, err := s.userCreator.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

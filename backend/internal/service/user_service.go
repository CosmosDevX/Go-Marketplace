package service

import (
	"context"
	"fmt"
	"myapp/internal/transport/dto"

	"golang.org/x/crypto/bcrypt"
)

type UserCreater interface {
	Create(ctx context.Context, userDTO dto.CreateUserDTO) (int, error)
}

type UserService struct {
	userCreater UserCreater
}

func NewUserService(userCreater UserCreater) UserService {
	return UserService{
		userCreater: userCreater,
	}
}

func (s UserService) Create(ctx context.Context, userDTO dto.CreateUserDTO) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		return 0, fmt.Errorf("password hashing: %w", err)
	}
	userDTO.Password = string(hashedPassword)

	id, err := s.userCreater.Create(ctx, userDTO)
	if err != nil {
		return 0, err
	}

	return id, nil
}

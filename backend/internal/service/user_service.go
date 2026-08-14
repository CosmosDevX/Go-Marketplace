package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/transport/dto"

	"golang.org/x/crypto/bcrypt"
)

type UserCreater interface {
	CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (int, *domain.DomainError)
}

type UserService struct {
	userCreater UserCreater
}

func NewUserService(userRepository UserCreater) UserService {
	return UserService{
		userCreater: userRepository,
	}
}

func (s UserService) CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (int, *domain.DomainError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		return 0, domain.NewDomainError(constants.InvalidPassword, "error during password hashing")
	}
	userDTO.Password = string(hashedPassword)

	id, apiErr := s.userCreater.CreateUser(ctx, userDTO)
	if apiErr != nil {
		return 0, apiErr
	}

	return id, nil
}

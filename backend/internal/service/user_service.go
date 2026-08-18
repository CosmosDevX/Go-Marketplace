package service

import (
	"context"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string
	Password string
	Email    string
}

type UserCreator interface {
	Create(ctx context.Context, u domain.User) (int, error)
}

type UserRoleCreator interface {
	Create(ctx context.Context, userID int, roleName string) error
}

type UserService struct {
	userCreator     UserCreator
	userRoleCreator UserRoleCreator
}

func NewUserService(userCreator UserCreator, userRoleCreator UserRoleCreator) UserService {
	return UserService{
		userCreator:     userCreator,
		userRoleCreator: userRoleCreator,
	}
}

func (s UserService) Create(ctx context.Context, input CreateUserInput) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return 0, fmt.Errorf("password hashing: %w", err)
	}

	passwordHash := string(hashedPassword)

	user, err := domain.NewUser(input.Username, passwordHash, input.Email)
	if err != nil {
		return 0, err
	}

	//TODO: add transaction: create user, create user role
	id, err := s.userCreator.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	if err := s.userRoleCreator.Create(ctx, id, config.DefaultUserRole); err != nil {
		return 0, err
	}

	return id, nil
}

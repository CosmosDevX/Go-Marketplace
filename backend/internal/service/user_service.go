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

type UserService struct {
	unitOfWork UnitOfWork
}

func NewUserService(unitOfWork UnitOfWork) UserService {
	return UserService{
		unitOfWork: unitOfWork,
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

	value, err := s.unitOfWork.Do(ctx, func(ctx context.Context, repos Repositories) (any, error) {
		id, err := repos.UserRepository.Create(ctx, user)
		if err != nil {
			return 0, err
		}

		if _, err := repos.CartRepository.Create(ctx, id); err != nil {
			return nil, err
		}

		if err := repos.UserRoleRepository.Create(ctx, id, config.DefaultUserRole); err != nil {
			return 0, err
		}

		return id, nil
	})

	if err != nil {
		return 0, err
	}

	id, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("id parse: %w", domain.ErrParse)
	}

	return id, nil
}

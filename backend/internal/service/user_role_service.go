package service

import (
	"context"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/domain"
)

type UserRoleInput struct {
	Username string
	Role     string
}

type UserRoleRepository interface {
	Create(ctx context.Context, username, roleName string) error
	Delete(ctx context.Context, username, roleName string) error
	ListByUsername(ctx context.Context, username string) ([]string, error)
}

type UsernameGetter interface {
	GetNameByID(ctx context.Context, userID int) (string, error)
}

type UserRoleService struct {
	userRoleRepository UserRoleRepository
	usernameGetter     UsernameGetter
}

func NewUserRoleService(userRoleRepository UserRoleRepository, usernameGetter UsernameGetter) UserRoleService {
	return UserRoleService{
		userRoleRepository: userRoleRepository,
		usernameGetter:     usernameGetter,
	}
}

func (s UserRoleService) Create(ctx context.Context, input UserRoleInput) error {
	if err := s.userRoleRepository.Create(ctx, input.Username, input.Role); err != nil {
		return err
	}

	return nil
}

func (s UserRoleService) Delete(ctx context.Context, input UserRoleInput, currentUserID int) error {
	username, err := s.usernameGetter.GetNameByID(ctx, currentUserID)
	if err != nil {
		return err
	}

	if username == input.Username && input.Role == config.AdminRole {
		return fmt.Errorf("cannot delete admin role from current user: %w", domain.ErrForbidden)
	}

	if err := s.userRoleRepository.Delete(ctx, input.Username, input.Role); err != nil {
		return err
	}

	return nil
}

func (s UserRoleService) List(ctx context.Context, username string) ([]string, error) {
	roles, err := s.userRoleRepository.ListByUsername(ctx, username)
	if err != nil {
		return []string{}, err
	}

	return roles, nil
}

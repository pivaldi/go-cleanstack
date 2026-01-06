package service

import (
	"context"
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/ports"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	logger := logging.GetLogger()

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	logger.Info("creating user", logging.String("email", user.Email))

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}

	return created, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from repository: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email from repository: %w", err)
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	users, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users from repository: %w", err)
	}

	return users, total, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	logger := logging.GetLogger()

	logger.Info("updating user", logging.Int64("id", user.ID))

	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user in repository: %w", err)
	}

	return updated, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	logger := logging.GetLogger()

	logger.Info("deleting user", logging.Int64("id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user from repository: %w", err)
	}

	return nil
}

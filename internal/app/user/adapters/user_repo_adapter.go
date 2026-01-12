package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/user/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/user/domain/ports"
	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence"
)

// UserRepositoryAdapter adapts the infra repository to the domain port.
// Since persistence uses entity.User directly, this is a simple pass-through
// that only translates error types.
type UserRepositoryAdapter struct {
	infraRepo *persistence.UserRepo
}

func NewUserRepositoryAdapter(infraRepo *persistence.UserRepo) ports.UserRepository {
	return &UserRepositoryAdapter{infraRepo: infraRepo}
}

// Ensure interface compliance.
var _ ports.UserRepository = (*UserRepositoryAdapter)(nil)

func (a *UserRepositoryAdapter) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	result, err := a.infraRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to create user: %w", err)
	}

	return result, nil
}

func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := a.infraRepo.GetByID(ctx, id)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to get user by id: %w", err)
	}

	return user, nil
}

func (a *UserRepositoryAdapter) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := a.infraRepo.GetByEmail(ctx, email)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to get user by email: %w", err)
	}

	return user, nil
}

func (a *UserRepositoryAdapter) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	users, total, err := a.infraRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("adapter: failed to list users: %w", err)
	}

	return users, total, nil
}

func (a *UserRepositoryAdapter) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	result, err := a.infraRepo.Update(ctx, user)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to update user: %w", err)
	}

	return result, nil
}

func (a *UserRepositoryAdapter) Delete(ctx context.Context, id int64) error {
	err := a.infraRepo.Delete(ctx, id)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return ports.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("adapter: failed to delete user: %w", err)
	}

	return nil
}

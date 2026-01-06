package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pivaldi/presence"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/ports"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/infra/persistence"
)

// UserRepositoryAdapter adapts the infra repository to the domain port.
// This is the bridge between domain (entities) and infra (DTOs).
type UserRepositoryAdapter struct {
	infraRepo *persistence.UserRepo
}

func NewUserRepositoryAdapter(infraRepo *persistence.UserRepo) ports.UserRepository {
	return &UserRepositoryAdapter{infraRepo: infraRepo}
}

// Ensure interface compliance.
var _ ports.UserRepository = (*UserRepositoryAdapter)(nil)

func (a *UserRepositoryAdapter) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dto := a.entityToDTO(user)

	resultDTO, err := a.infraRepo.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to create user: %w", err)
	}

	return a.dtoToEntity(resultDTO), nil
}

func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	dto, err := a.infraRepo.GetByID(ctx, id)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to get user by id: %w", err)
	}

	return a.dtoToEntity(dto), nil
}

func (a *UserRepositoryAdapter) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dto, err := a.infraRepo.GetByEmail(ctx, email)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to get user by email: %w", err)
	}

	return a.dtoToEntity(dto), nil
}

func (a *UserRepositoryAdapter) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	dtos, total, err := a.infraRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("adapter: failed to list users: %w", err)
	}

	users := make([]*entity.User, len(dtos))
	for i, dto := range dtos {
		users[i] = a.dtoToEntity(dto)
	}

	return users, total, nil
}

func (a *UserRepositoryAdapter) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	dto := a.entityToDTO(user)

	resultDTO, err := a.infraRepo.Update(ctx, dto)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return nil, ports.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("adapter: failed to update user: %w", err)
	}

	return a.dtoToEntity(resultDTO), nil
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

func (a *UserRepositoryAdapter) entityToDTO(user *entity.User) *persistence.UserDTO {
	dto := &persistence.UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role.String(),
	}

	if user.FirstName.IsSet() && !user.FirstName.IsNull() {
		dto.FirstName = sql.NullString{String: user.FirstName.MustGet(), Valid: true}
	}

	if user.LastName.IsSet() && !user.LastName.IsNull() {
		dto.LastName = sql.NullString{String: user.LastName.MustGet(), Valid: true}
	}

	return dto
}

func (a *UserRepositoryAdapter) dtoToEntity(dto *persistence.UserDTO) *entity.User {
	role, _ := entity.ParseRole(dto.Role)

	user := &entity.User{
		ID:        dto.ID,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      role,
		CreatedAt: dto.CreatedAt,
	}

	if dto.FirstName.Valid {
		user.FirstName = presence.FromValue(dto.FirstName.String)
	}

	if dto.LastName.Valid {
		user.LastName = presence.FromValue(dto.LastName.String)
	}

	if dto.UpdatedAt.Valid {
		user.UpdatedAt = presence.FromValue(dto.UpdatedAt.Time)
	}

	if dto.DeletedAt.Valid {
		user.DeletedAt = presence.FromValue(dto.DeletedAt.Time)
	}

	return user
}

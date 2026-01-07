package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pivaldi/go-cleanstack/internal/app/user/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/user/domain/ports"
	"github.com/pivaldi/go-cleanstack/internal/app/user/service"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logger/zap"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize no-op logger for all tests
	logging.SetLogger(zap.NewNop())
}

// MockUserRepository is a mock implementation of ports.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(
	ctx context.Context,
	offset, limit int,
) ([]*entity.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Ensure MockUserRepository implements ports.UserRepository.
var _ ports.UserRepository = (*MockUserRepository)(nil)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := entity.NewUser("test@example.com", "password123", entity.RoleUser)
		input.SetFirstName("John")
		input.SetLastName("Doe")

		expected := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "hashed_password",
			FirstName: presence.FromValue("John"),
			LastName:  presence.FromValue("Doe"),
			Role:      entity.RoleUser,
			CreatedAt: time.Now(),
		}

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

		result, err := svc.CreateUser(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - empty email", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := entity.NewUser("", "password123", entity.RoleUser)

		result, err := svc.CreateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation failed")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("validation error - invalid email format", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := entity.NewUser("invalid-email", "password123", entity.RoleUser)

		result, err := svc.CreateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation failed")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("validation error - password too short", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := entity.NewUser("test@example.com", "short", entity.RoleUser)

		result, err := svc.CreateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation failed")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := entity.NewUser("test@example.com", "password123", entity.RoleUser)
		repoErr := errors.New("database connection failed")

		mockRepo.On("Create", mock.Anything, input).Return(nil, repoErr)

		result, err := svc.CreateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create user in repository")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		expected := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "hashed_password",
			FirstName: presence.FromValue("John"),
			LastName:  presence.FromValue("Doe"),
			Role:      entity.RoleUser,
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := svc.GetUserByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, ports.ErrUserNotFound)

		result, err := svc.GetUserByID(context.Background(), 999)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get user from repository")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		repoErr := errors.New("database error")
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, repoErr)

		result, err := svc.GetUserByID(context.Background(), 1)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get user from repository")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		expected := &entity.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "hashed_password",
			Role:      entity.RoleAdmin,
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expected, nil)

		result, err := svc.GetUserByEmail(context.Background(), "test@example.com")

		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Email, result.Email)
		assert.Equal(t, entity.RoleAdmin, result.Role)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, ports.ErrUserNotFound)

		result, err := svc.GetUserByEmail(context.Background(), "notfound@example.com")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get user by email from repository")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_ListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		users := []*entity.User{
			{ID: 1, Email: "user1@example.com", Role: entity.RoleUser},
			{ID: 2, Email: "user2@example.com", Role: entity.RoleAdmin},
		}

		mockRepo.On("List", mock.Anything, 0, 10).Return(users, int64(2), nil)

		result, total, err := svc.ListUsers(context.Background(), 0, 10)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		mockRepo.On("List", mock.Anything, 0, 10).Return([]*entity.User{}, int64(0), nil)

		result, total, err := svc.ListUsers(context.Background(), 0, 10)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("pagination", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		users := []*entity.User{
			{ID: 11, Email: "user11@example.com", Role: entity.RoleUser},
			{ID: 12, Email: "user12@example.com", Role: entity.RoleUser},
		}

		mockRepo.On("List", mock.Anything, 10, 5).Return(users, int64(15), nil)

		result, total, err := svc.ListUsers(context.Background(), 10, 5)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(15), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		repoErr := errors.New("database error")
		mockRepo.On("List", mock.Anything, 0, 10).Return(nil, int64(0), repoErr)

		result, total, err := svc.ListUsers(context.Background(), 0, 10)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "failed to list users from repository")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := &entity.User{
			ID:        1,
			Email:     "updated@example.com",
			FirstName: presence.FromValue("Jane"),
			LastName:  presence.FromValue("Smith"),
			Role:      entity.RoleAdmin,
		}

		expected := &entity.User{
			ID:        1,
			Email:     "updated@example.com",
			Password:  "hashed_password",
			FirstName: presence.FromValue("Jane"),
			LastName:  presence.FromValue("Smith"),
			Role:      entity.RoleAdmin,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: presence.FromValue(time.Now()),
		}

		mockRepo.On("Update", mock.Anything, input).Return(expected, nil)

		result, err := svc.UpdateUser(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Email, result.Email)
		assert.True(t, result.UpdatedAt.IsSet())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := &entity.User{
			ID:    999,
			Email: "notfound@example.com",
			Role:  entity.RoleUser,
		}

		mockRepo.On("Update", mock.Anything, input).Return(nil, ports.ErrUserNotFound)

		result, err := svc.UpdateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to update user in repository")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		input := &entity.User{
			ID:    1,
			Email: "test@example.com",
			Role:  entity.RoleUser,
		}
		repoErr := errors.New("database error")

		mockRepo.On("Update", mock.Anything, input).Return(nil, repoErr)

		result, err := svc.UpdateUser(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to update user in repository")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := svc.DeleteUser(context.Background(), 1)

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(999)).Return(ports.ErrUserNotFound)

		err := svc.DeleteUser(context.Background(), 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete user from repository")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := service.NewUserService(mockRepo)

		repoErr := errors.New("database error")
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(repoErr)

		err := svc.DeleteUser(context.Background(), 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete user from repository")
		mockRepo.AssertExpectations(t)
	})
}

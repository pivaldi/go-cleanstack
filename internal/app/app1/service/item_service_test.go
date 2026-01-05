package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/entity"
)

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(ctx context.Context, item *entity.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id string) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemRepository) List(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestItemService_CreateItem_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	svc := NewItemService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(nil)

	item, err := svc.CreateItem(context.Background(), "Test Item", "Test Description")

	assert.NoError(t, err)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "Test Description", item.Description)
	mockRepo.AssertExpectations(t)
}

func TestItemService_CreateItem_ValidationError(t *testing.T) {
	mockRepo := new(MockItemRepository)
	svc := NewItemService(mockRepo)

	item, err := svc.CreateItem(context.Background(), "", "Description")

	assert.Error(t, err)
	assert.Nil(t, item)
	mockRepo.AssertNotCalled(t, "Create")
}

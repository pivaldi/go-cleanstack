package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/pivaldi/go-cleanstack/internal/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/domain/ports"
)

type ItemService struct {
	repo ports.ItemRepository
}

func NewItemService(repo ports.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) CreateItem(ctx context.Context, name, description string) (*entity.Item, error) {
	item := entity.NewItem(uuid.New().String(), name, description)

	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("item validation failed: %w", err)
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create item in repository: %w", err)
	}

	return item, nil
}

func (s *ItemService) GetItem(ctx context.Context, id string) (*entity.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from repository: %w", err)
	}

	return item, nil
}

func (s *ItemService) ListItems(ctx context.Context) ([]*entity.Item, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list items from repository: %w", err)
	}

	return items, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item from repository: %w", err)
	}

	return nil
}

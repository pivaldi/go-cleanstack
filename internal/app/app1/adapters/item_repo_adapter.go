package adapters

import (
	"context"
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/ports"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/infra/persistence"
)

// ItemRepositoryAdapter adapts the infra repository to the domain port
// This is the bridge between domain (entities) and infra (DTOs)
type ItemRepositoryAdapter struct {
	infraRepo *persistence.ItemRepo
}

func NewItemRepositoryAdapter(infraRepo *persistence.ItemRepo) ports.ItemRepository {
	return &ItemRepositoryAdapter{infraRepo: infraRepo}
}

// Ensure interface compliance
var _ ports.ItemRepository = (*ItemRepositoryAdapter)(nil)

func (a *ItemRepositoryAdapter) Create(ctx context.Context, item *entity.Item) error {
	dto := &persistence.ItemDTO{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
	}

	if err := a.infraRepo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (a *ItemRepositoryAdapter) GetByID(ctx context.Context, id string) (*entity.Item, error) {
	dto, err := a.infraRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item by id: %w", err)
	}

	return &entity.Item{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		CreatedAt:   dto.CreatedAt,
	}, nil
}

func (a *ItemRepositoryAdapter) List(ctx context.Context) ([]*entity.Item, error) {
	dtos, err := a.infraRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	items := make([]*entity.Item, len(dtos))
	for i, dto := range dtos {
		items[i] = &entity.Item{
			ID:          dto.ID,
			Name:        dto.Name,
			Description: dto.Description,
			CreatedAt:   dto.CreatedAt,
		}
	}

	return items, nil
}

func (a *ItemRepositoryAdapter) Delete(ctx context.Context, id string) error {
	if err := a.infraRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

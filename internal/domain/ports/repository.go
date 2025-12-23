package ports

import (
	"context"

	"github.com/pivaldi/go-cleanstack/internal/domain/entity"
)

type ItemRepository interface {
	Create(ctx context.Context, item *entity.Item) error
	GetByID(ctx context.Context, id string) (*entity.Item, error)
	List(ctx context.Context) ([]*entity.Item, error)
	Delete(ctx context.Context, id string) error
}

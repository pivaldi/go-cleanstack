//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pivaldi/go-cleanstack/internal/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence"
	"github.com/pivaldi/go-cleanstack/tests/testutil"
)

func TestItemRepo_CRUD(t *testing.T) {
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	// Setup test DB
	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := persistence.NewItemRepo(db, logger)

	t.Run("Create and GetByID", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item := entity.NewItem("test-id-1", "Test Item", "Test Description")

		// Convert to DTO
		dto := &persistence.ItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}

		err := repo.Create(ctx, dto)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, "test-id-1")
		require.NoError(t, err)
		assert.Equal(t, item.ID, retrieved.ID)
		assert.Equal(t, item.Name, retrieved.Name)
		assert.Equal(t, item.Description, retrieved.Description)
	})

	t.Run("List", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item1 := entity.NewItem("test-id-2", "Item 1", "Desc 1")
		item2 := entity.NewItem("test-id-3", "Item 2", "Desc 2")

		dto1 := &persistence.ItemDTO{
			ID:          item1.ID,
			Name:        item1.Name,
			Description: item1.Description,
			CreatedAt:   item1.CreatedAt,
		}
		dto2 := &persistence.ItemDTO{
			ID:          item2.ID,
			Name:        item2.Name,
			Description: item2.Description,
			CreatedAt:   item2.CreatedAt,
		}

		require.NoError(t, repo.Create(ctx, dto1))
		require.NoError(t, repo.Create(ctx, dto2))

		items, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item := entity.NewItem("test-id-4", "To Delete", "Will be deleted")
		dto := &persistence.ItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}
		require.NoError(t, repo.Create(ctx, dto))

		err := repo.Delete(ctx, "test-id-4")
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, "test-id-4")
		assert.Error(t, err)
	})
}

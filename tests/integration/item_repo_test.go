//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pivaldi/go-cleanstack/internal/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence"
	"github.com/pivaldi/go-cleanstack/internal/platform/logging"
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

	logger := logging.NewNop()
	repo := persistence.NewItemRepo(db, logger)

	t.Run("Create and GetByID", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		id := uuid.New().String()
		item := entity.NewItem(id, "Test Item", "Test Description")

		// Convert to DTO
		dto := &persistence.ItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}

		err := repo.Create(ctx, dto)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, item.ID, retrieved.ID)
		assert.Equal(t, item.Name, retrieved.Name)
		assert.Equal(t, item.Description, retrieved.Description)
	})

	t.Run("List", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		id1 := uuid.New().String()
		id2 := uuid.New().String()
		item1 := entity.NewItem(id1, "Item 1", "Desc 1")
		item2 := entity.NewItem(id2, "Item 2", "Desc 2")

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

		id := uuid.New().String()
		item := entity.NewItem(id, "To Delete", "Will be deleted")
		dto := &persistence.ItemDTO{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}
		require.NoError(t, repo.Create(ctx, dto))

		err := repo.Delete(ctx, id)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, id)
		assert.Error(t, err)
	})
}

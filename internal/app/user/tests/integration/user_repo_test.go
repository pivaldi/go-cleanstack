//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pivaldi/go-cleanstack/internal/app/user/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence"
	"github.com/pivaldi/go-cleanstack/internal/app/user/tests/testutil"
)

func TestUserRepo_CRUD(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	repo := persistence.NewUserRepo(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("test@example.com", "password123", entity.RoleUser)

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "test@example.com", created.Email)
		assert.NotEqual(t, "password123", created.Password) // Password should be hashed

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Email, retrieved.Email)
	})

	t.Run("Create with optional fields", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("full@example.com", "password123", entity.RoleAdmin)
		user.SetFirstName("John")
		user.SetLastName("Doe")

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.True(t, created.FirstName.IsSet())
		assert.False(t, created.FirstName.IsNull())
		assert.Equal(t, "John", created.FirstName.MustGet())
		assert.True(t, created.LastName.IsSet())
		assert.False(t, created.LastName.IsNull())
		assert.Equal(t, "Doe", created.LastName.MustGet())
		assert.Equal(t, entity.RoleAdmin, created.Role)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("email@example.com", "password123", entity.RoleUser)

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByEmail(ctx, "email@example.com")
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
	})

	t.Run("List with pagination", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create 5 users
		for i := 0; i < 5; i++ {
			user := entity.NewUser(
				"user"+string(rune('0'+i))+"@example.com",
				"password123",
				entity.RoleUser,
			)

			_, err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		// Get first page
		users, total, err := repo.List(ctx, 0, 2)
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(5), total)

		// Get second page
		users, total, err = repo.List(ctx, 2, 2)
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(5), total)

		// Get last page
		users, total, err = repo.List(ctx, 4, 2)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, int64(5), total)
	})

	t.Run("Update", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("update@example.com", "password123", entity.RoleUser)

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update user
		created.Email = "updated@example.com"
		created.SetFirstName("Jane")
		created.Role = entity.RoleAdmin

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", updated.Email)
		assert.Equal(t, "Jane", updated.FirstName.MustGet())
		assert.Equal(t, entity.RoleAdmin, updated.Role)
		assert.True(t, updated.UpdatedAt.IsSet())
	})

	t.Run("Update password", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("passupdate@example.com", "oldpassword1", entity.RoleUser)

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		oldHash := created.Password

		// Update password
		created.Password = "newpassword1"

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.NotEqual(t, oldHash, updated.Password)        // Password hash should change
		assert.NotEqual(t, "newpassword1", updated.Password) // Should be hashed
	})

	t.Run("Delete (soft delete)", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		user := entity.NewUser("delete@example.com", "password123", entity.RoleUser)

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Should not be found
		_, err = repo.GetByID(ctx, created.ID)
		assert.ErrorIs(t, err, persistence.ErrUserNotFound)

		// Should not be in list
		users, total, err := repo.List(ctx, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, users)
		assert.Equal(t, int64(0), total)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		_, err := repo.GetByID(ctx, 9999)
		assert.ErrorIs(t, err, persistence.ErrUserNotFound)
	})

	t.Run("GetByEmail not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		_, err := repo.GetByEmail(ctx, "notfound@example.com")
		assert.ErrorIs(t, err, persistence.ErrUserNotFound)
	})

	t.Run("Delete not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		err := repo.Delete(ctx, 9999)
		assert.ErrorIs(t, err, persistence.ErrUserNotFound)
	})
}

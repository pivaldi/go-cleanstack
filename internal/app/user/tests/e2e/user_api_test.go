//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pivaldi/go-cleanstack/internal/app/user/adapters"
	userv1 "github.com/pivaldi/go-cleanstack/internal/app/user/api/gen/user/v1"
	"github.com/pivaldi/go-cleanstack/internal/app/user/api/gen/user/v1/userv1connect"
	"github.com/pivaldi/go-cleanstack/internal/app/user/api/handler"
	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence"
	"github.com/pivaldi/go-cleanstack/internal/app/user/service"
	"github.com/pivaldi/go-cleanstack/internal/app/user/tests/testutil"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logger/zap"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
)

var l logging.Logger

func init() {
	var err error
	l, err = zap.NewDevelopment("debug")
	if err != nil {
		panic(err)
	}
}

func TestUserAPI_E2E(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	// Wire up dependencies
	infraRepo := persistence.NewUserRepo(db)
	userRepo := adapters.NewUserRepositoryAdapter(infraRepo)
	userService := service.NewUserService(userRepo, l)
	userHandler := handler.NewUserHandler(userService)

	// Create test server
	mux := http.NewServeMux()
	path, h := userv1connect.NewUserServiceHandler(userHandler)
	mux.Handle(path, h)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create client
	client := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		server.URL,
	)

	t.Run("CreateUser", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		firstName := "John"
		lastName := "Doe"
		resp, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:     "e2e@example.com",
			Password:  "password123",
			FirstName: &firstName,
			LastName:  &lastName,
			Role:      "user",
		}))

		require.NoError(t, err)
		assert.NotZero(t, resp.Msg.User.Id)
		assert.Equal(t, "e2e@example.com", resp.Msg.User.Email)
		assert.Equal(t, "John", *resp.Msg.User.FirstName)
		assert.Equal(t, "Doe", *resp.Msg.User.LastName)
		assert.Equal(t, "user", resp.Msg.User.Role)
	})

	t.Run("CreateUser with validation error", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		_, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:    "invalid-email",
			Password: "password123",
			Role:     "user",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("GetUser", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create user first
		createResp, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:    "get@example.com",
			Password: "password123",
			Role:     "user",
		}))
		require.NoError(t, err)

		// Get user
		getResp, err := client.GetUser(ctx, connect.NewRequest(&userv1.GetUserRequest{
			Id: createResp.Msg.User.Id,
		}))
		require.NoError(t, err)
		assert.Equal(t, createResp.Msg.User.Id, getResp.Msg.User.Id)
		assert.Equal(t, "get@example.com", getResp.Msg.User.Email)
	})

	t.Run("GetUser not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		_, err := client.GetUser(ctx, connect.NewRequest(&userv1.GetUserRequest{
			Id: 9999,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create user first
		_, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:    "byemail@example.com",
			Password: "password123",
			Role:     "admin",
		}))
		require.NoError(t, err)

		// Get by email
		getResp, err := client.GetUserByEmail(ctx, connect.NewRequest(&userv1.GetUserByEmailRequest{
			Email: "byemail@example.com",
		}))
		require.NoError(t, err)
		assert.Equal(t, "byemail@example.com", getResp.Msg.User.Email)
		assert.Equal(t, "admin", getResp.Msg.User.Role)
	})

	t.Run("ListUsers", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create users
		for i := 0; i < 3; i++ {
			_, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
				Email:    fmt.Sprintf("list%d@example.com", i),
				Password: "password123",
				Role:     "user",
			}))
			require.NoError(t, err)
		}

		// List users
		listResp, err := client.ListUsers(ctx, connect.NewRequest(&userv1.ListUsersRequest{
			Offset: 0,
			Limit:  10,
		}))
		require.NoError(t, err)
		assert.Len(t, listResp.Msg.Users, 3)
		assert.Equal(t, int64(3), listResp.Msg.Total)
	})

	t.Run("ListUsers with pagination", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create users
		for i := 0; i < 5; i++ {
			_, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
				Email:    fmt.Sprintf("page%d@example.com", i),
				Password: "password123",
				Role:     "user",
			}))
			require.NoError(t, err)
		}

		// First page
		listResp, err := client.ListUsers(ctx, connect.NewRequest(&userv1.ListUsersRequest{
			Offset: 0,
			Limit:  2,
		}))
		require.NoError(t, err)
		assert.Len(t, listResp.Msg.Users, 2)
		assert.Equal(t, int64(5), listResp.Msg.Total)

		// Second page
		listResp, err = client.ListUsers(ctx, connect.NewRequest(&userv1.ListUsersRequest{
			Offset: 2,
			Limit:  2,
		}))
		require.NoError(t, err)
		assert.Len(t, listResp.Msg.Users, 2)
		assert.Equal(t, int64(5), listResp.Msg.Total)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create user
		createResp, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:    "update@example.com",
			Password: "password123",
			Role:     "user",
		}))
		require.NoError(t, err)

		// Update user
		newEmail := "updated@example.com"
		newFirstName := "Updated"
		updateResp, err := client.UpdateUser(ctx, connect.NewRequest(&userv1.UpdateUserRequest{
			Id:        createResp.Msg.User.Id,
			Email:     &newEmail,
			FirstName: &newFirstName,
		}))
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", updateResp.Msg.User.Email)
		assert.Equal(t, "Updated", *updateResp.Msg.User.FirstName)
		assert.NotNil(t, updateResp.Msg.User.UpdatedAt)
	})

	t.Run("UpdateUser not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		newEmail := "updated@example.com"
		_, err := client.UpdateUser(ctx, connect.NewRequest(&userv1.UpdateUserRequest{
			Id:    9999,
			Email: &newEmail,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	})

	t.Run("DeleteUser", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create user
		createResp, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
			Email:    "delete@example.com",
			Password: "password123",
			Role:     "user",
		}))
		require.NoError(t, err)

		// Delete user
		_, err = client.DeleteUser(ctx, connect.NewRequest(&userv1.DeleteUserRequest{
			Id: createResp.Msg.User.Id,
		}))
		require.NoError(t, err)

		// Verify deleted
		_, err = client.GetUser(ctx, connect.NewRequest(&userv1.GetUserRequest{
			Id: createResp.Msg.User.Id,
		}))
		assert.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	})

	t.Run("DeleteUser not found", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		_, err := client.DeleteUser(ctx, connect.NewRequest(&userv1.DeleteUserRequest{
			Id: 9999,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	})
}

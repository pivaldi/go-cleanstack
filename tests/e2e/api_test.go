//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pivaldi/go-cleanstack/internal/app/adapters"
	"github.com/pivaldi/go-cleanstack/internal/app/service"
	cleanstackv1 "github.com/pivaldi/go-cleanstack/internal/infra/api/gen/cleanstack/v1"
	"github.com/pivaldi/go-cleanstack/internal/infra/api/gen/cleanstack/v1/cleanstackv1connect"
	"github.com/pivaldi/go-cleanstack/internal/infra/api/handler"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence"
	"github.com/pivaldi/go-cleanstack/tests/testutil"
)

func TestItemAPI_E2E(t *testing.T) {
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	// Setup test DB
	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	// Wire up dependencies
	logger := zap.NewNop()
	infraRepo := persistence.NewItemRepo(db, logger)
	itemRepo := adapters.NewItemRepositoryAdapter(infraRepo)
	itemService := service.NewItemService(itemRepo, logger)
	itemHandler := handler.NewItemHandler(itemService, logger)

	// Create test server
	mux := http.NewServeMux()
	path, h := cleanstackv1connect.NewItemServiceHandler(itemHandler)
	mux.Handle(path, h)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create client
	client := cleanstackv1connect.NewItemServiceClient(
		http.DefaultClient,
		server.URL,
	)

	t.Run("CreateItem", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		resp, err := client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "E2E Test Item",
			Description: "Created via E2E test",
		}))

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Msg.Item.Id)
		assert.Equal(t, "E2E Test Item", resp.Msg.Item.Name)
	})

	t.Run("ListItems", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create items first
		_, err := client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "Item 1",
			Description: "First item",
		}))
		require.NoError(t, err)

		_, err = client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "Item 2",
			Description: "Second item",
		}))
		require.NoError(t, err)

		// List items
		resp, err := client.ListItems(ctx, connect.NewRequest(&cleanstackv1.ListItemsRequest{}))
		require.NoError(t, err)
		assert.Len(t, resp.Msg.Items, 2)
	})
}

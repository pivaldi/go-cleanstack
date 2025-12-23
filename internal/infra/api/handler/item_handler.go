package handler

import (
	"context"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/pivaldi/go-cleanstack/internal/app/service"
	cleanstackv1 "github.com/pivaldi/go-cleanstack/internal/infra/api/gen/cleanstack/v1"
	"github.com/pivaldi/go-cleanstack/internal/infra/api/gen/cleanstack/v1/cleanstackv1connect"
)

type ItemHandler struct {
	service *service.ItemService
	logger  *zap.Logger
}

func NewItemHandler(svc *service.ItemService, logger *zap.Logger) *ItemHandler {
	return &ItemHandler{
		service: svc,
		logger:  logger,
	}
}

var _ cleanstackv1connect.ItemServiceHandler = (*ItemHandler)(nil)

func (h *ItemHandler) CreateItem(
	ctx context.Context,
	req *connect.Request[cleanstackv1.CreateItemRequest],
) (*connect.Response[cleanstackv1.CreateItemResponse], error) {
	item, err := h.service.CreateItem(ctx, req.Msg.Name, req.Msg.Description)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&cleanstackv1.CreateItemResponse{
		Item: &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}), nil
}

func (h *ItemHandler) GetItem(
	ctx context.Context,
	req *connect.Request[cleanstackv1.GetItemRequest],
) (*connect.Response[cleanstackv1.GetItemResponse], error) {
	item, err := h.service.GetItem(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&cleanstackv1.GetItemResponse{
		Item: &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}), nil
}

func (h *ItemHandler) ListItems(
	ctx context.Context,
	_ *connect.Request[cleanstackv1.ListItemsRequest],
) (*connect.Response[cleanstackv1.ListItemsResponse], error) {
	items, err := h.service.ListItems(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoItems := make([]*cleanstackv1.Item, len(items))
	for i, item := range items {
		protoItems[i] = &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return connect.NewResponse(&cleanstackv1.ListItemsResponse{
		Items: protoItems,
	}), nil
}

func (h *ItemHandler) DeleteItem(
	ctx context.Context,
	req *connect.Request[cleanstackv1.DeleteItemRequest],
) (*connect.Response[cleanstackv1.DeleteItemResponse], error) {
	if err := h.service.DeleteItem(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cleanstackv1.DeleteItemResponse{}), nil
}

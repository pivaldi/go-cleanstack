package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/entity"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/domain/ports"
	app1v1 "github.com/pivaldi/go-cleanstack/internal/app/app1/infra/api/gen/app1/v1"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/infra/api/gen/app1/v1/app1v1connect"
	"github.com/pivaldi/go-cleanstack/internal/app/app1/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

var _ app1v1connect.UserServiceHandler = (*UserHandler)(nil)

func (h *UserHandler) CreateUser(
	ctx context.Context,
	req *connect.Request[app1v1.CreateUserRequest],
) (*connect.Response[app1v1.CreateUserResponse], error) {
	role, err := entity.ParseRole(req.Msg.Role)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	user := entity.NewUser(req.Msg.Email, req.Msg.Password, role)

	if req.Msg.FirstName != nil {
		user.SetFirstName(*req.Msg.FirstName)
	}
	if req.Msg.LastName != nil {
		user.SetLastName(*req.Msg.LastName)
	}

	created, err := h.service.CreateUser(ctx, user)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&app1v1.CreateUserResponse{
		User: h.entityToProto(created),
	}), nil
}

func (h *UserHandler) GetUser(
	ctx context.Context,
	req *connect.Request[app1v1.GetUserRequest],
) (*connect.Response[app1v1.GetUserResponse], error) {
	user, err := h.service.GetUserByID(ctx, req.Msg.Id)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&app1v1.GetUserResponse{
		User: h.entityToProto(user),
	}), nil
}

func (h *UserHandler) GetUserByEmail(
	ctx context.Context,
	req *connect.Request[app1v1.GetUserByEmailRequest],
) (*connect.Response[app1v1.GetUserByEmailResponse], error) {
	user, err := h.service.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&app1v1.GetUserByEmailResponse{
		User: h.entityToProto(user),
	}), nil
}

func (h *UserHandler) ListUsers(
	ctx context.Context,
	req *connect.Request[app1v1.ListUsersRequest],
) (*connect.Response[app1v1.ListUsersResponse], error) {
	users, total, err := h.service.ListUsers(ctx, int(req.Msg.Offset), int(req.Msg.Limit))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoUsers := make([]*app1v1.User, len(users))
	for i, user := range users {
		protoUsers[i] = h.entityToProto(user)
	}

	return connect.NewResponse(&app1v1.ListUsersResponse{
		Users: protoUsers,
		Total: total,
	}), nil
}

func (h *UserHandler) UpdateUser(
	ctx context.Context,
	req *connect.Request[app1v1.UpdateUserRequest],
) (*connect.Response[app1v1.UpdateUserResponse], error) {
	user := &entity.User{ID: req.Msg.Id}

	if req.Msg.Email != nil {
		user.Email = *req.Msg.Email
	}
	if req.Msg.Password != nil {
		user.Password = *req.Msg.Password
	}
	if req.Msg.FirstName != nil {
		user.SetFirstName(*req.Msg.FirstName)
	}
	if req.Msg.LastName != nil {
		user.SetLastName(*req.Msg.LastName)
	}
	if req.Msg.Role != nil {
		role, err := entity.ParseRole(*req.Msg.Role)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		user.Role = role
	}

	updated, err := h.service.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&app1v1.UpdateUserResponse{
		User: h.entityToProto(updated),
	}), nil
}

func (h *UserHandler) DeleteUser(
	ctx context.Context,
	req *connect.Request[app1v1.DeleteUserRequest],
) (*connect.Response[app1v1.DeleteUserResponse], error) {
	if err := h.service.DeleteUser(ctx, req.Msg.Id); err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&app1v1.DeleteUserResponse{}), nil
}

func (h *UserHandler) entityToProto(user *entity.User) *app1v1.User {
	proto := &app1v1.User{
		Id:        user.ID,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.FirstName.IsSet() && !user.FirstName.IsNull() {
		val := user.FirstName.MustGet()
		proto.FirstName = &val
	}

	if user.LastName.IsSet() && !user.LastName.IsNull() {
		val := user.LastName.MustGet()
		proto.LastName = &val
	}

	if user.UpdatedAt.IsSet() && !user.UpdatedAt.IsNull() {
		val := user.UpdatedAt.MustGet()
		formatted := val.Format("2006-01-02T15:04:05Z07:00")
		proto.UpdatedAt = &formatted
	}

	return proto
}

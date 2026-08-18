package grpchandler

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/sixsat/user-manager-be-challenge/port"
	userproto "github.com/sixsat/user-manager-be-challenge/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type handler struct {
	userproto.UnimplementedUserServiceServer
	validate *validator.Validate
	userSvc  port.UserService
}

func New(validate *validator.Validate, userSvc port.UserService) *handler {
	return &handler{
		validate: validate,
		userSvc:  userSvc,
	}
}

var _ userproto.UserServiceServer = (*handler)(nil)

func (h *handler) CreateUser(ctx context.Context, req *userproto.CreateUserReq) (*userproto.CreateUserRes, error) {
	dto := createUserReq{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	if err := h.validate.Struct(&dto); err != nil {
		slog.Error("[grpc handler] error validating request", slog.String("error", err.Error()))
		return nil, invalidRequestError(err)
	}

	if err := h.userSvc.Create(ctx, dto.toDomain()); err != nil {
		return nil, serviceError(err)
	}

	return &userproto.CreateUserRes{}, nil
}

func (h *handler) GetUser(ctx context.Context, req *userproto.GetUserReq) (*userproto.GetUserRes, error) {
	dto := getUserReq{ID: req.GetId()}
	if err := h.validate.Struct(&dto); err != nil {
		slog.Error("[grpc handler] error validating request", slog.String("error", err.Error()))
		return nil, invalidRequestError(err)
	}

	user, err := h.userSvc.GetByID(ctx, dto.ID)
	if err != nil {
		return nil, serviceError(err)
	}

	return &userproto.GetUserRes{
		User: &userproto.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

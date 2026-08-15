package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) port.UserRepository {
	return &userRepo{
		collection: client.Database("user_manager").Collection("users"),
	}
}

func (r *userRepo) Create(ctx context.Context, req *domain.CreateUserReq) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrDuplicateUser
		}
		slog.Error("[mongodb] error inserting user", slog.String("error", err.Error()))
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.GetUserRes, error) {
	return nil, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.GetUserByEmailRes, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		slog.Error("[mongodb] error finding user by email", slog.String("error", err.Error()))
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &domain.GetUserByEmailRes{
		ID:           user.ID.Hex(),
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *userRepo) List(ctx context.Context) ([]domain.GetUserRes, error) {
	return nil, nil
}

func (r *userRepo) Update(ctx context.Context, req *domain.UpdateUserReq) error {
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *userRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}

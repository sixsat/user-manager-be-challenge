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
		CreatedAt:    time.Now(),
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
	objectID, _ := bson.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		slog.Error("[mongodb] error finding user by id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &domain.GetUserRes{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		slog.Error("[mongodb] error finding users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("find users: %w", err)
	}

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		slog.Error("[mongodb] error decoding users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("decode users: %w", err)
	}

	res := make([]domain.GetUserRes, 0, len(users))
	for _, u := range users {
		res = append(res, domain.GetUserRes{
			ID:        u.ID.Hex(),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}
	return res, nil
}

func (r *userRepo) Update(ctx context.Context, req *domain.UpdateUserReq) error {
	objectID, _ := bson.ObjectIDFromHex(req.ID)

	fields := bson.M{}
	if req.Name != nil {
		fields["name"] = *req.Name
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": fields})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrDuplicateUser
		}
		slog.Error("[mongodb] error updating user", slog.String("error", err.Error()))
		return fmt.Errorf("update user: %w", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	objectID, _ := bson.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		slog.Error("[mongodb] error deleting user", slog.String("error", err.Error()))
		return fmt.Errorf("delete user: %w", err)
	}
	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepo) Count(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		slog.Error("[mongodb] error counting users", slog.String("error", err.Error()))
		return 0, fmt.Errorf("count users: %w", err)
	}
	return int(count), nil
}

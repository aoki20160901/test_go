// internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"myapi/internal/model"
	"strconv"
	"time"
)

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	// DBアクセス（例）
	return &model.User{
		ID:   id,
		Name: "Taro",
	}, nil
}

// CreateUser creates a new user (mock implementation)
func CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}
	// Simple ID generation for example purposes
	u.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	return u, nil
}

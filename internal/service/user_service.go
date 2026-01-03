// internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"myapi/internal/model"
	"myapi/internal/repository"
)

func GetUser(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return repository.FindUserByID(ctx, id)
}

// CreateUser creates a new user via the repository
func CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}
	return repository.CreateUser(ctx, u)
}

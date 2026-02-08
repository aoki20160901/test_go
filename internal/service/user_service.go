// internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"myapi/internal/model"
	"myapi/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(
	userRepo repository.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	id string,
) (*model.User, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.userRepo.FindByID(ctx, id)
}

// CreateUser creates a new user via the repository
func (s *UserService) CreateUser(
	ctx context.Context,
	u *model.User,
) (*model.User, error) {

	if u == nil {
		return nil, errors.New("user is nil")
	}

	return s.userRepo.Create(ctx, u)
}

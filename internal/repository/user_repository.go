package repository

import (
	"context"
	"errors"
	"myapi/internal/model"

	"gorm.io/gorm"
)

// -----------------------------
// Repository interface
// -----------------------------
type UserRepository interface {
	Create(ctx context.Context, u *model.User) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

// -----------------------------
// GORM 実装
// -----------------------------
type userRepository struct {
	db *gorm.DB
}

// コンストラクタ
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}
	if r.db == nil {
		return nil, errors.New("db is nil")
	}
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if r.db == nil {
		return nil, errors.New("db is nil")
	}
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

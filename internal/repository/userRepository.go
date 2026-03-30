package repository

import (
	"context"
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserById(ctx context.Context, userId uint) (*domain.User, error)
	GetUserByUuid(ctx context.Context, userUuid uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("get user by id %d: %w", id, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByUuid(ctx context.Context, userUuid uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("uuid = ?", userUuid).First(&user).Error; err != nil {
		return nil, fmt.Errorf("get user by uuid %s: %w", userUuid, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Updates() with a struct skips zero-value fields (false, 0, "").
	// This is intentional here — UpdateUser is used for partial updates.
	// Use Save() instead if you ever need to explicitly clear a field to its zero value.
	err := r.db.WithContext(ctx).Model(user).Clauses(clause.Returning{}).Where("uuid = ?", user.Uuid).Updates(user).Error
	if err != nil {
		return fmt.Errorf("update user %s: %w", user.Uuid, err)
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, user *domain.User) error {
	// GORM soft-deletes because User has gorm.DeletedAt.
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return fmt.Errorf("delete user %s: %w", user.Uuid, err)
	}
	return nil
}

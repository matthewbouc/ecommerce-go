package repository

import (
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUserById(userId uint) (*domain.User, error)
	GetUserByUuid(userUuid uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetUserById(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("get user by id %d: %w", id, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByUuid(userUuid uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("uuid = ?", userUuid).First(&user).Error; err != nil {
		return nil, fmt.Errorf("get user by uuid %s: %w", userUuid, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {

	err := r.db.Model(user).Clauses(clause.Returning{}).Where("uuid = ?", user.Uuid).Updates(user).Error
	if err != nil {
		return fmt.Errorf("update user %s: %w", user.Uuid, err)
	}
	return nil
}

func (r *userRepository) DeleteUser(user *domain.User) error {
	// GORM soft-deletes because User has gorm.DeletedAt.
	if err := r.db.Delete(user).Error; err != nil {
		return fmt.Errorf("delete user %s: %w", user.Uuid, err)
	}
	return nil
}

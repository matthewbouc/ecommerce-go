package repository

import (
	"ecommerce/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserById(userId uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	//UpdateUser(user domain.User) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *domain.User) error {

	err := r.db.Create(&user).Error
	if err != nil {
		return fmt.Errorf("create user error: %w", err)

	}
	return nil
}

func (r *userRepository) GetUserById(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("get user by id %d: %w", id, err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}
	return &user, nil
}

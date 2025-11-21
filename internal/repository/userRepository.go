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
	// gorm handles created_at and updated_at
	err := r.db.Create(&user).Error
	if err != nil {
		return domain.User{}, fmt.Errorf("create user error: %w", err)

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
	if err := r.db.First(&user, "uuid = ?", userUuid).Error; err != nil {
		return nil, fmt.Errorf("error during get user by uuid %s: %w", userUuid, err)
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

func (r *userRepository) UpdateUser(user *domain.User) error {
	// gorm will update_at
	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", user.Id).Updates(user).Error
	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}
	return nil
}

func (r *userRepository) DeleteUser(user *domain.User) error {
	// gorm will soft-delete
	err := r.db.Delete(&user).Error
	if err != nil {
		return fmt.Errorf("error while deleting user: %w", err)
	}
	return nil
}

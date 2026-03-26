package repository

import (
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserId(userId uint) ([]domain.CartItem, error)
	GetCartItemByUserAndProduct(userId uint, productId uuid.UUID) (*domain.CartItem, error)
	CreateCartItem(item domain.CartItem) (domain.CartItem, error)
	UpdateCartItem(item *domain.CartItem) error
	DeleteCartItem(id uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserId(userId uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	if err := r.db.Where("user_id = ?", userId).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get cart for user %d: %w", userId, err)
	}
	return items, nil
}

func (r *cartRepository) GetCartItemByUserAndProduct(userId uint, productId uuid.UUID) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userId, productId).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) CreateCartItem(item domain.CartItem) (domain.CartItem, error) {
	if err := r.db.Create(&item).Error; err != nil {
		return domain.CartItem{}, fmt.Errorf("create cart item: %w", err)
	}
	return item, nil
}

func (r *cartRepository) UpdateCartItem(item *domain.CartItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return fmt.Errorf("update cart item %d: %w", item.Id, err)
	}
	return nil
}

func (r *cartRepository) DeleteCartItem(id uint) error {
	if err := r.db.Delete(&domain.CartItem{}, id).Error; err != nil {
		return fmt.Errorf("delete cart item %d: %w", id, err)
	}
	return nil
}

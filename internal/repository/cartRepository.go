package repository

import (
	"context"
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserId(ctx context.Context, userId uint) ([]domain.CartItem, error)
	GetCartItemByUserAndProduct(ctx context.Context, userId uint, productId uuid.UUID) (*domain.CartItem, error)
	CreateCartItem(ctx context.Context, item domain.CartItem) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, item *domain.CartItem) error
	DeleteCartItemByIdAndUser(ctx context.Context, itemId uint, userId uint) error
	ClearCartByUserId(ctx context.Context, userId uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserId(ctx context.Context, userId uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get cart for user %d: %w", userId, err)
	}
	return items, nil
}

func (r *cartRepository) GetCartItemByUserAndProduct(ctx context.Context, userId uint, productId uuid.UUID) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userId, productId).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) CreateCartItem(ctx context.Context, item domain.CartItem) (domain.CartItem, error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return domain.CartItem{}, fmt.Errorf("create cart item: %w", err)
	}
	return item, nil
}

func (r *cartRepository) UpdateCartItem(ctx context.Context, item *domain.CartItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update cart item %d: %w", item.Id, err)
	}
	return nil
}

func (r *cartRepository) ClearCartByUserId(ctx context.Context, userId uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.CartItem{}).Error; err != nil {
		return fmt.Errorf("clear cart for user %d: %w", userId, err)
	}
	return nil
}

func (r *cartRepository) DeleteCartItemByIdAndUser(ctx context.Context, itemId uint, userId uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", itemId, userId).Delete(&domain.CartItem{})
	if result.Error != nil {
		return fmt.Errorf("delete cart item %d: %w", itemId, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("cart item not found")
	}
	return nil
}

package repository

import (
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrdersByUserId(userId uint) ([]domain.Order, error)
	GetOrderByUuid(userId uint, orderUuid uuid.UUID) (*domain.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetOrdersByUserId(userId uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("get orders for user %d: %w", userId, err)
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByUuid(userId uint, orderUuid uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Where("user_id = ? AND uuid = ?", userId, orderUuid).First(&order).Error
	if err != nil {
		return nil, fmt.Errorf("get order %s: %w", orderUuid, err)
	}
	return &order, nil
}

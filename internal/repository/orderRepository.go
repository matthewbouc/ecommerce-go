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
	CreateOrder(order domain.Order) (domain.Order, error)
	UpdateOrderStatus(orderUuid uuid.UUID, status domain.OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order domain.Order) (domain.Order, error) {
	if err := r.db.Create(&order).Error; err != nil {
		return domain.Order{}, fmt.Errorf("create order: %w", err)
	}
	return order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderUuid uuid.UUID, status domain.OrderStatus) error {
	result := r.db.Model(&domain.Order{}).
		Where("uuid = ?", orderUuid).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update order status %s: %w", orderUuid, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order %s not found", orderUuid)
	}
	return nil
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

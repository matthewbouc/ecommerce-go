package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

type Order struct {
	Id         uint        `json:"id" gorm:"column:id;primaryKey"`
	Uuid       uuid.UUID   `json:"uuid" gorm:"column:uuid;type:uuid"`
	UserId     uint        `json:"user_id" gorm:"column:user_id;not null"`
	Status     OrderStatus `json:"status" gorm:"column:status;default:pending"`
	TotalPrice float64     `json:"total_price" gorm:"column:total_price"`
	CreatedAt  time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderId"`
	User       User        `json:"-" gorm:"foreignKey:UserId"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.Uuid == uuid.Nil {
		o.Uuid = uuid.New()
	}
	return nil
}

func (o *Order) Validate() error {
	if o.UserId == 0 {
		return errors.New("user ID is required")
	}
	if o.TotalPrice < 0 {
		return errors.New("total price cannot be negative")
	}
	if o.Status != "" && !o.Status.IsValidStatus() {
		return errors.New("invalid order status")
	}
	return nil
}

func (s OrderStatus) IsValidStatus() bool {
	switch s {
	case OrderPending, OrderProcessing, OrderShipped, OrderDelivered, OrderCancelled:
		return true
	}
	return false
}

type OrderItem struct {
	Id        uint      `json:"id" gorm:"column:id;primaryKey"`
	OrderId   uint      `json:"order_id" gorm:"column:order_id;not null"`
	ProductId uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	SellerId  uuid.UUID `json:"seller_id" gorm:"column:seller_id;type:uuid;not null"`
	Name      string    `json:"name" gorm:"column:name"`
	ImageUrl  string    `json:"image_url" gorm:"column:image_url"`
	Price     float64   `json:"price" gorm:"column:price"`
	Quantity  uint      `json:"quantity" gorm:"column:quantity"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (oi *OrderItem) Validate() error {
	if oi.OrderId == 0 {
		return errors.New("order ID is required")
	}
	if oi.ProductId == uuid.Nil {
		return errors.New("product ID is required")
	}
	if oi.SellerId == uuid.Nil {
		return errors.New("seller ID is required")
	}
	if oi.Name == "" {
		return errors.New("product name is required")
	}
	if oi.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if oi.Quantity == 0 {
		return errors.New("quantity must be at least 1")
	}
	return nil
}

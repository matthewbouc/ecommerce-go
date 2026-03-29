package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	Id        uint      `json:"id" gorm:"column:id;primaryKey"`
	UserId    uint      `json:"user_id" gorm:"column:user_id;not null"`
	ProductId uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	SellerId  uuid.UUID `json:"seller_id" gorm:"column:seller_id;type:uuid;not null"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	ImageUrl  string    `json:"image_url" gorm:"column:image_url"`
	// price is snapshotted at add-time for display; validate Catalog price at checkout.
	Price     float64   `json:"price" gorm:"column:price;not null"`
	Quantity  uint      `json:"quantity" gorm:"column:quantity;not null;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	User      User      `json:"-" gorm:"foreignKey:UserId"`
}

func (c *CartItem) Validate() error {
	if c.UserId == 0 {
		return errors.New("user ID is required")
	}
	if c.ProductId == uuid.Nil {
		return errors.New("product ID is required")
	}
	if c.SellerId == uuid.Nil {
		return errors.New("seller ID is required")
	}
	if c.Name == "" {
		return errors.New("product name is required")
	}
	if c.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if c.Quantity == 0 {
		return errors.New("quantity must be at least 1")
	}
	return nil
}

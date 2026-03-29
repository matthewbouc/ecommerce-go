package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Id          uint           `json:"id" gorm:"column:id;primaryKey"`
	Uuid        uuid.UUID      `json:"uuid" gorm:"column:uuid;type:uuid;uniqueIndex;not null"`
	SellerID    uuid.UUID      `json:"seller_id" gorm:"column:seller_id;type:uuid;index;not null"`
	Name        string         `json:"name" gorm:"column:name;not null"`
	Description string         `json:"description" gorm:"column:description"`
	Price       float64        `json:"price" gorm:"column:price;type:numeric(10,2);not null"`
	Stock       int            `json:"stock" gorm:"column:stock;not null;default:0"`
	CategoryID  uuid.UUID      `json:"category_id" gorm:"column:category_id;type:uuid;index;not null"`
	Category    Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:Uuid"`
	ImageUrl    string         `json:"image_url" gorm:"column:image_url"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.Uuid == uuid.Nil {
		p.Uuid = uuid.New()
	}
	return nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	if p.SellerID == uuid.Nil {
		return errors.New("seller ID is required")
	}
	if p.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	return nil
}

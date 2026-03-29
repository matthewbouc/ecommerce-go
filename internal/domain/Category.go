package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	Id          uint      `json:"id" gorm:"column:id;primaryKey"`
	Uuid        uuid.UUID `json:"uuid" gorm:"column:uuid;type:uuid;uniqueIndex;not null"`
	Name        string    `json:"name" gorm:"column:name;uniqueIndex;not null"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Uuid == uuid.Nil {
		c.Uuid = uuid.New()
	}
	return nil
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	return nil
}

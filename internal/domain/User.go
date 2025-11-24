package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uint           `json:"id" gorm:"column:id;primaryKey"`
	Uuid             uuid.UUID      `json:"uuid" gorm:"column:uuid;type:uuid"`
	FirstName        *string        `json:"first_name" gorm:"column:first_name"`
	LastName         *string        `json:"last_name" gorm:"column:last_name"`
	Phone            *string        `json:"phone" gorm:"column:phone"`
	Email            string         `json:"email" gorm:"column:email;index;unique;not null"`
	Password         string         `json:"password" gorm:"column:password;not null"`
	VerificationCode int            `json:"verification_code" gorm:"column:verification_code"`
	Expiry           time.Time      `json:"expiry" gorm:"column:expiry;default:null"`
	Verified         bool           `json:"verified" gorm:"column:verified;default:false"`
	UserType         string         `json:"user_type" gorm:"column:user_type;default:buyer"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	LastLogin        *time.Time     `json:"last_login" gorm:"column:last_login;default:null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Uuid == uuid.Nil {
		u.Uuid = uuid.New()
	}
	return nil
}

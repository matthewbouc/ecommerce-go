package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id               uint      `json:"id" gorm:"primaryKey"`
	Uuid             uuid.UUID `json:"uuid"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Phone            *string   `json:"phone"`
	Email            string    `json:"email" gorm:"index;unique;not null"`
	Password         string    `json:"password"`
	VerificationCode string    `json:"verification_code"`
	Expiry           time.Time `json:"expiry"`
	Verified         bool      `json:"verified" gorm:"default:false"`
	UserType         string    `json:"user_type" gorm:"default:buyer"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"default:now()"`
}

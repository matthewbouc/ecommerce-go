package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id               uint      `json:"id"`
	Uuid             uuid.UUID `json:"uuid"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Password         string    `json:"password"`
	VerificationCode string    `json:"verification_code"`
	Expiry           time.Time `json:"expiry"`
	Verified         bool      `json:"verified"`
	UserType         string    `json:"user_type"`
}

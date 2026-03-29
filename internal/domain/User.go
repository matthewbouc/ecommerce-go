package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserType string

const (
	SELLER = "seller"
	BUYER  = "buyer"
)

type User struct {
	Id               uint           `json:"id" gorm:"column:id;primaryKey"`
	Uuid             uuid.UUID      `json:"uuid" gorm:"column:uuid;type:uuid"`
	FirstName        string         `json:"first_name" gorm:"column:first_name"`
	LastName         string         `json:"last_name" gorm:"column:last_name"`
	Phone            string         `json:"phone" gorm:"column:phone"`
	Email            string         `json:"email" gorm:"column:email;index;unique;not null"`
	Password         string         `json:"password" gorm:"column:password;not null"`
	VerificationCode int            `json:"verification_code" gorm:"column:verification_code"`
	Expiry           *time.Time     `json:"expiry" gorm:"column:expiry;default:null"`
	Verified         bool           `json:"verified" gorm:"column:verified;default:false"`
	UserType         UserType       `json:"user_type" gorm:"column:user_type;default:buyer"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	LastLogin        *time.Time     `json:"last_login" gorm:"column:last_login;default:null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Uuid == uuid.Nil {
		u.Uuid = uuid.New()
	}

	if u.UserType == "" {
		u.UserType = BUYER
	}
	if !u.UserType.IsValidUserType() {
		return errors.New("invalid user type")
	}
	return nil
}

func (ut UserType) IsValidUserType() bool {
	return ut == BUYER || ut == SELLER
}

func (u *User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	// Basic sanity check — a proper email regex is overkill at the domain layer
	if !containsAt(u.Email) {
		return errors.New("email is invalid")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.Phone == "" {
		return errors.New("phone number is required")
	}
	if u.UserType != "" && !u.UserType.IsValidUserType() {
		return errors.New("user type must be 'buyer' or 'seller'")
	}
	return nil
}

func containsAt(email string) bool {
	for _, c := range email {
		if c == '@' {
			return true
		}
	}
	return false
}

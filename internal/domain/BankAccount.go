package domain

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	Id                uint           `json:"id" gorm:"primary_key;auto_increment"`
	UserId            uint           `json:"user_id" gorm:"column:user_id;not null"`
	User              User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BankAccountNumber uint           `json:"bank_account_number" gorm:"index;column:bank_account_number;not null"`
	SwiftCode         string         `json:"swift_code" gorm:"column:swift_code;not null"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}

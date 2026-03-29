package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	Id                uint           `json:"id" gorm:"primary_key;auto_increment"`
	UserId            uint           `json:"user_id" gorm:"column:user_id;not null"`
	User              User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BankAccountNumber uint           `json:"bank_account_number" gorm:"index;column:bank_account_number;not null"`
	RoutingNumber     uint           `json:"routing_number" gorm:"column:routing_number;null"`
	SwiftCode         uint           `json:"swift_code" gorm:"column:swift_code;default:null"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}

func (b *BankAccount) Validate() error {
	if b.UserId == 0 {
		return errors.New("user ID is required")
	}
	if b.BankAccountNumber == 0 {
		return errors.New("bank account number is required")
	}
	// RoutingNumber is required for domestic ACH transfers.
	// SwiftCode is required for international wires.
	// Require at least one to be present.
	if b.RoutingNumber == 0 && b.SwiftCode == 0 {
		return errors.New("routing number or swift code is required")
	}
	return nil
}

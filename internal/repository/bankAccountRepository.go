package repository

import (
	"ecommerce/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

type BankAccountRepository interface {
	CreateBankAccount(bankAccount domain.BankAccount) (domain.BankAccount, error)
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepository{
		db: db,
	}
}

func (r *bankAccountRepository) CreateBankAccount(bankAccount domain.BankAccount) (domain.BankAccount, error) {
	err := r.db.Create(&bankAccount).Error
	if err != nil {
		return domain.BankAccount{}, fmt.Errorf("create user error: %w", err)
	}
	return bankAccount, nil
}

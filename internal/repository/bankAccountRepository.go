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
	if err := r.db.Create(&bankAccount).Error; err != nil {
		return domain.BankAccount{}, fmt.Errorf("create bank account: %w", err)
	}
	return bankAccount, nil
}

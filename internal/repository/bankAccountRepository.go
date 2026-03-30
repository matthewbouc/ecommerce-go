package repository

import (
	"context"
	"ecommerce/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

type BankAccountRepository interface {
	CreateBankAccount(ctx context.Context, bankAccount domain.BankAccount) (domain.BankAccount, error)
}

type bankAccountRepository struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepository{db: db}
}

func (r *bankAccountRepository) CreateBankAccount(ctx context.Context, bankAccount domain.BankAccount) (domain.BankAccount, error) {
	if err := r.db.WithContext(ctx).Create(&bankAccount).Error; err != nil {
		return domain.BankAccount{}, fmt.Errorf("create bank account: %w", err)
	}
	return bankAccount, nil
}

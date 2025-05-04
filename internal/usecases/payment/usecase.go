package paymentcase

import (
	"context"
	"fmt"
	"myproject/internal/entities"
	paymentrepo "myproject/internal/repositories/payment"
)

type PaymentUseCase interface {
	Deposit(ctx context.Context, userID int, amount float64) error
	CreateTransaction(ctx context.Context, tx *entities.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID int) ([]entities.Transaction, error)
}

type service struct {
	repo paymentrepo.Repository
}

func NewPaymentService(repo paymentrepo.Repository) PaymentUseCase {
	return &service{repo: repo}
}

func (s *service) Deposit(ctx context.Context, userID int, amount float64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}
	if amount <= 0 {
		return fmt.Errorf("invalid deposit amount: %.2f", amount)
	}
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *service) CreateTransaction(ctx context.Context, tx *entities.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction data cannot be nil")
	}
	if tx.UserID <= 0 {
		return fmt.Errorf("invalid user ID in transaction: %d", tx.UserID)
	}
	if tx.Amount == 0 {
		return fmt.Errorf("transaction amount cannot be zero")
	}
	if tx.Type == "" {
		return fmt.Errorf("transaction type cannot be empty")
	}
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *service) GetTransactionsByUser(ctx context.Context, userID int) ([]entities.Transaction, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	return s.repo.GetTransactionsByUserID(ctx, userID)
}

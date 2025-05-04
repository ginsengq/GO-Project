package paymentservice

import (
	"context"
	"fmt"
	"myproject/internal/entities"
	paymentrepo "myproject/internal/repositories/payment"
)

type PaymentService interface {
	Deposit(ctx context.Context, userID int, amount float64) error
	CreateTransaction(ctx context.Context, tx *entities.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID int) ([]entities.Transaction, error)
}

type Service struct {
	repo paymentrepo.Repository
}

func NewService(repo paymentrepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Deposit(ctx context.Context, userID int, amount float64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}
	if amount <= 0 {
		return fmt.Errorf("invalid deposit amount: %.2f", amount)
	}
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *Service) CreateTransaction(ctx context.Context, tx *entities.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction cannot be nil")
	}
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *Service) GetTransactionsByUser(ctx context.Context, userID int) ([]entities.Transaction, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	return s.repo.GetTransactionsByUserID(ctx, userID)
}

package paymentrepo

import (
	"context"
	"myproject/internal/entities"
)

type Repository interface {
	ProcessPayment(ctx context.Context, payment *entities.Payment) error
	GetPaymentByID(ctx context.Context, paymentID int) (*entities.Payment, error)
	Deposit(ctx context.Context, userID int, amount float64) error // Add this method
	CreateTransaction(ctx context.Context, tx *entities.Transaction) error
	GetTransactionsByUserID(ctx context.Context, userID int) ([]entities.Transaction, error)
}

package paymentrepo

import (
	"context"
	"fmt"
	"myproject/internal/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) ProcessPayment(ctx context.Context, payment *entities.Payment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, payment.Amount, payment.UserID)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	transaction := &entities.Transaction{
		UserID:      payment.UserID,
		Amount:      payment.Amount,
		Type:        payment.PaymentMethod,
		Description: fmt.Sprintf("Payment via %s", payment.PaymentMethod),
		CreatedAt:   payment.CreatedAt,
	}
	err = r.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE payments SET status = $1, transaction_id = $2 WHERE id = $3`, payment.Status, transaction.ID, payment.ID)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *repository) Deposit(ctx context.Context, userID int, amount float64) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to deposit: %w", err)
	}
	return nil
}

func (r *repository) GetPaymentByID(ctx context.Context, paymentID int) (*entities.Payment, error) {
	row := r.db.QueryRow(ctx, "SELECT id, user_id, amount, payment_method, status, transaction_id, created_at, provider_id FROM payments WHERE id = $1", paymentID)

	var payment entities.Payment
	err := row.Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.CreatedAt,
		&payment.ProviderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &payment, nil
}

func (r *repository) CreateTransaction(ctx context.Context, tx *entities.Transaction) error {
	query := `INSERT INTO transactions (user_id, amount, type, description, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, tx.UserID, tx.Amount, tx.Type, tx.Description, tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

func (r *repository) GetTransactionsByUserID(ctx context.Context, userID int) ([]entities.Transaction, error) {
	query := `SELECT id, user_id, amount, type, description, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var tx entities.Transaction
		err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Description, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

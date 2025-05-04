package orderrepo

import (
	"context"
	"myproject/internal/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, order *entities.Order) (int, error) {
	query := `
		INSERT INTO orders (user_id, car_id, status, deposit, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	var id int
	err := r.db.QueryRow(ctx, query, order.UserID, order.CarID, order.Status, order.Deposit, order.TotalPrice).Scan(&id)
	return id, err
}

func (r *repository) GetByID(ctx context.Context, id int) (*entities.Order, error) {
	query := `SELECT id, user_id, car_id, status, deposit, total_price, created_at, updated_at FROM orders WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var order entities.Order
	err := row.Scan(&order.ID, &order.UserID, &order.CarID, &order.Status, &order.Deposit, &order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID int) ([]entities.Order, error) {
	query := `SELECT id, user_id, car_id, status, deposit, total_price, created_at, updated_at FROM orders WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	for rows.Next() {
		var order entities.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.CarID, &order.Status, &order.Deposit, &order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = current_timestamp WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListAll(ctx context.Context) ([]entities.Order, error) {
	query := `SELECT id, user_id, car_id, status, deposit, total_price, created_at, updated_at FROM orders`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	for rows.Next() {
		var order entities.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.CarID, &order.Status, &order.Deposit, &order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

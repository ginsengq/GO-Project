package userrepo

import (
	"context"
	"errors"
	entity "myproject/internal/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (name, email, password_hash, balance, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.PasswordHash, user.Balance, user.Role)
	return err
}

func (r *postgresRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, name, email, password_hash, balance, role FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Balance, &user.Role,
	)
	if err != nil {
		if errors.Is(err, errors.New("pgx: no rows in result")) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, name, email, password_hash, balance, role FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Balance, &user.Role,
	)
	if err != nil {
		if errors.Is(err, errors.New("pgx: no rows in result")) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name=$1, email=$2, balance=$3, role=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Balance, user.Role, user.ID)
	return err
}

func (r *postgresRepo) UpdateBalance(ctx context.Context, id int, amount float64) error {
	query := `UPDATE users SET balance=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, amount, id)
	return err
}

func (r *postgresRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresRepo) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *postgresRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	query := `SELECT id, name, email, password_hash, balance, role FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Balance, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *postgresRepo) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

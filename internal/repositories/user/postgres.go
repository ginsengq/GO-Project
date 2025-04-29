package user

import (
	"context"
	"myproject/internal/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(user entity.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, user.Name, user.Email, user.PasswordHash).Scan(&id)
	return id, err
}

func (r *postgresRepo) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, name, email, password_hash, balance, role FROM users WHERE email = $1`
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Balance, &user.Role,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) GetByID(id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, name, email, password_hash, balance, role FROM users WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Balance, &user.Role,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) Update(user entity.User) error {
	query := `UPDATE users SET name=$1, email=$2, balance=$3, role=$4 WHERE id=$5`
	_, err := r.db.Exec(context.Background(), query, user.Name, user.Email, user.Balance, user.Role, user.ID)
	return err
}

func (r *postgresRepo) Delete(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

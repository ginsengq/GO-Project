package repositories

import (
	"myproject/internal/repositories/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	UserRepo user.Repository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepo: user.NewPostgresRepo(db),
	}
}

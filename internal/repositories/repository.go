package repositories

import (
	carrepo "myproject/internal/repositories/car"
	orderrepo "myproject/internal/repositories/order"
	paymentrepo "myproject/internal/repositories/payment"
	userrepo "myproject/internal/repositories/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	User    userrepo.Repository
	Car     carrepo.Repository
	Order   orderrepo.Repository
	Payment paymentrepo.Repository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		User:    userrepo.NewPostgresRepo(db),
		Car:     carrepo.NewPostgresRepo(db),
		Order:   orderrepo.NewPostgresRepository(db),
		Payment: paymentrepo.NewPaymentRepository(db),
	}
}

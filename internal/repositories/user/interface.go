package userrepo

import (
	"context"
	"myproject/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id int) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	UpdateBalance(ctx context.Context, id int, amount float64) error
	Delete(ctx context.Context, id int) error
	IsEmailExists(ctx context.Context, email string) (bool, error)
	List(ctx context.Context, limit, offset int) ([]*entities.User, error) // Добавьте этот метод
	Count(ctx context.Context) (int, error)                                // Добавьте этот метод
}

package orderrepo

import (
	"context"
	"myproject/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, order *entities.Order) (int, error)
	GetByID(ctx context.Context, id int) (*entities.Order, error)
	GetByUserID(ctx context.Context, userID int) ([]entities.Order, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]entities.Order, error)
}

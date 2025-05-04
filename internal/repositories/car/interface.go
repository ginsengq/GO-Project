package carrepo

import (
	"context"
	"myproject/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, car *entities.Car) (int, error)
	GetByID(ctx context.Context, id int) (*entities.Car, error)
	Update(ctx context.Context, id int, update entities.CarUpdate) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, error)
	SetStatus(ctx context.Context, id int, status string) error
}

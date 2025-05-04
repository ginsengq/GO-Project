package car

import (
	"context"
	"fmt"
	"myproject/internal/entities"
	carrepo "myproject/internal/repositories/car"
)

type CarUseCase interface {
	CreateCar(ctx context.Context, input *entities.Car) (*entities.Car, error)
	GetCar(ctx context.Context, id int) (*entities.Car, error)
	UpdateCar(ctx context.Context, id int, input entities.CarUpdate) (*entities.Car, error)
	DeleteCar(ctx context.Context, id int) error
	ListCars(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, int, error)
	ChangeCarStatus(ctx context.Context, id int, status entities.CarStatus) (*entities.Car, error)
}

type carUseCase struct {
	repo carrepo.Repository
}

func NewCarUseCase(repo carrepo.Repository) CarUseCase {
	return &carUseCase{repo: repo}
}

func (uc *carUseCase) CreateCar(ctx context.Context, input *entities.Car) (*entities.Car, error) {
	if input.Brand == "" || input.Model == "" || input.Year < 1900 {
		return nil, fmt.Errorf("invalid input: brand, model, and year are required")
	}

	input.Status = entities.CarStatusAvailable
	id, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create car: %w", err)
	}

	createdCar, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created car: %w", err)
	}
	return createdCar, nil
}

func (uc *carUseCase) GetCar(ctx context.Context, id int) (*entities.Car, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid input: car ID must be a positive integer")
	}

	car, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get car by ID: %w", err)
	}
	return car, nil
}

func (uc *carUseCase) UpdateCar(ctx context.Context, id int, input entities.CarUpdate) (*entities.Car, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid input: car ID must be a positive integer")
	}

	err := uc.repo.Update(ctx, id, input)
	if err != nil {
		return nil, fmt.Errorf("update car: %w", err)
	}

	updatedCar, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated car: %w", err)
	}
	return updatedCar, nil
}

func (uc *carUseCase) DeleteCar(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid input: car ID must be a positive integer")
	}

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete car: %w", err)
	}
	return nil
}

func (uc *carUseCase) ListCars(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, int, error) {
	cars, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list cars: %w", err)
	}
	return cars, len(cars), nil
}

func (uc *carUseCase) ChangeCarStatus(ctx context.Context, id int, status entities.CarStatus) (*entities.Car, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid input: car ID must be a positive integer")
	}

	if status != entities.CarStatusAvailable && status != entities.CarStatusReserved && status != entities.CarStatusSold {
		return nil, fmt.Errorf("invalid car status: %s", status)
	}

	err := uc.repo.SetStatus(ctx, id, string(status))
	if err != nil {
		return nil, fmt.Errorf("change car status: %w", err)
	}

	updatedCar, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated car: %w", err)
	}
	return updatedCar, nil
}

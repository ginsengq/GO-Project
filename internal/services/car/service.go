package carservice

import (
	"context"
	"errors"
	"fmt"
	"myproject/internal/entities"
	carrepo "myproject/internal/repositories/car"
	"time"
)

type CarService interface {
	CreateCar(ctx context.Context, input *entities.Car) (*entities.Car, error)
	GetCar(ctx context.Context, id int) (*entities.Car, error)
	UpdateCar(ctx context.Context, id int, input entities.CarUpdate) (*entities.Car, error)
	DeleteCar(ctx context.Context, id int) error
	ListCars(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, int, error)
	ChangeCarStatus(ctx context.Context, id int, status entities.CarStatus) (*entities.Car, error)
	CheckAvailability(ctx context.Context, carID int, startDate, endDate time.Time) (bool, error)
	UpdateStatus(ctx context.Context, carID int, status string) error
}

type service struct {
	repo carrepo.Repository
}

func NewService(repo carrepo.Repository) CarService {
	return &service{repo: repo}
}

func (s *service) CreateCar(ctx context.Context, input *entities.Car) (*entities.Car, error) {
	if err := validateCar(input); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	input.Status = entities.CarStatusAvailable
	id, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) GetCar(ctx context.Context, id int) (*entities.Car, error) {
	if id <= 0 {
		return nil, errors.New("invalid car ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateCar(ctx context.Context, id int, input entities.CarUpdate) (*entities.Car, error) {
	if id <= 0 {
		return nil, errors.New("invalid car ID")
	}

	if err := s.repo.Update(ctx, id, input); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) DeleteCar(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid car ID")
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) ListCars(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, int, error) {
	cars, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return cars, len(cars), nil
}

func (s *service) ChangeCarStatus(ctx context.Context, id int, status entities.CarStatus) (*entities.Car, error) {
	if id <= 0 {
		return nil, errors.New("invalid car ID")
	}

	if err := s.repo.SetStatus(ctx, id, string(status)); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) CheckAvailability(ctx context.Context, carID int, startDate, endDate time.Time) (bool, error) {
	car, err := s.repo.GetByID(ctx, carID)
	if err != nil {
		return false, fmt.Errorf("failed to get car: %w", err)
	}
	return car.Status == entities.CarStatusAvailable, nil
}

func (s *service) UpdateStatus(ctx context.Context, carID int, status string) error {
	if carID <= 0 {
		return errors.New("invalid car ID")
	}
	return s.repo.SetStatus(ctx, carID, status)
}

func validateCar(car *entities.Car) error {
	if car.Brand == "" {
		return errors.New("brand is required")
	}
	if car.Model == "" {
		return errors.New("model is required")
	}
	if car.Year < 1900 || car.Year > 2100 {
		return errors.New("invalid year")
	}
	if car.Price <= 0 {
		return errors.New("price must be positive")
	}
	if car.Mileage < 0 {
		return errors.New("mileage cannot be negative")
	}
	return nil
}

package ordercase

import (
	"context"
	"fmt"
	"myproject/internal/entities"
	orderrepo "myproject/internal/repositories/order"
)

type UseCase interface {
	CreateOrder(ctx context.Context, order *entities.Order) (int, error)
	GetOrder(ctx context.Context, id int) (*entities.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]entities.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) error
	CancelOrder(ctx context.Context, id int) error
	ListAllOrders(ctx context.Context) ([]entities.Order, error)
}

type service struct {
	repo orderrepo.Repository
}

func NewService(repo orderrepo.Repository) UseCase {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, order *entities.Order) (int, error) {
	if order == nil {
		return 0, fmt.Errorf("order cannot be nil")
	}
	if order.UserID <= 0 {
		return 0, fmt.Errorf("invalid user ID: %d", order.UserID)
	}
	if order.CarID <= 0 {
		return 0, fmt.Errorf("invalid car ID: %d", order.CarID)
	}
	if order.Deposit < 0 {
		return 0, fmt.Errorf("deposit cannot be negative: %.2f", order.Deposit)
	}
	if order.TotalPrice <= 0 {
		return 0, fmt.Errorf("total price must be positive: %.2f", order.TotalPrice)
	}
	order.Status = entities.OrderStatusPending // Установка статуса по умолчанию
	return s.repo.Create(ctx, order)
}

func (s *service) GetOrder(ctx context.Context, id int) (*entities.Order, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid order ID: %d", id)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetOrdersByUserID(ctx context.Context, userID int) ([]entities.Order, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return fmt.Errorf("invalid order ID: %d", id)
	}
	if status == "" {
		return fmt.Errorf("order status cannot be empty")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *service) CancelOrder(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid order ID: %d", id)
	}
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order with ID %d: %w", id, err)
	}
	if order.Status == entities.OrderStatusCancelled || order.Status == entities.OrderStatusCompleted {
		return fmt.Errorf("order with ID %d is already %s and cannot be cancelled", id, order.Status)
	}
	return s.repo.UpdateStatus(ctx, id, entities.OrderStatusCancelled)
}

func (s *service) ListAllOrders(ctx context.Context) ([]entities.Order, error) {
	return s.repo.ListAll(ctx)
}

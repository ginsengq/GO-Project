package orderservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"myproject/internal/entities"
	orderrepo "myproject/internal/repositories/order"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderData   = errors.New("invalid order data")
	ErrOrderAlreadyClosed = errors.New("order is already completed or cancelled")
	ErrInvalidStatus      = errors.New("invalid order status")
)

type Service struct {
	repo           orderrepo.Repository
	carService     CarService
	userService    UserService
	paymentService PaymentService
}

type CarService interface {
	CheckAvailability(ctx context.Context, carID int, startDate, endDate time.Time) (bool, error)
	UpdateStatus(ctx context.Context, carID int, status string) error
}

type UserService interface {
	CheckBalance(ctx context.Context, userID int, amount float64) (bool, error)
	DeductBalance(ctx context.Context, userID int, amount float64) error
}

type PaymentService interface {
	CreateTransaction(ctx context.Context, tx *entities.Transaction) error // Change the signature to take *entities.Transaction
}

func NewService(
	repo orderrepo.Repository,
	carService CarService,
	userService UserService,
	paymentService PaymentService,
) *Service {
	return &Service{
		repo:           repo,
		carService:     carService,
		userService:    userService,
		paymentService: paymentService,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *entities.Order) (int, error) {
	if err := validateOrder(order); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidOrderData, err)
	}

	available, err := s.carService.CheckAvailability(
		ctx,
		order.CarID,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to check car availability: %w", err)
	}
	if !available {
		return 0, errors.New("car is not available for selected dates")
	}

	hasBalance, err := s.userService.CheckBalance(ctx, order.UserID, order.TotalPrice)
	if err != nil {
		return 0, fmt.Errorf("failed to check user balance: %w", err)
	}
	if !hasBalance {
		return 0, errors.New("insufficient funds")
	}

	order.Status = entities.OrderStatusPending
	order.CreatedAt = time.Now()

	id, err := s.repo.Create(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.userService.DeductBalance(ctx, order.UserID, order.TotalPrice); err != nil {
		return 0, fmt.Errorf("failed to deduct balance: %w", err)
	}

	if err := s.carService.UpdateStatus(ctx, order.CarID, "reserved"); err != nil {
		return 0, fmt.Errorf("failed to update car status: %w", err)
	}

	transaction := &entities.Transaction{
		UserID:      order.UserID,
		Amount:      order.TotalPrice,
		Type:        "order_payment",
		Description: fmt.Sprintf("Payment for order #%d", id),
		CreatedAt:   time.Now(),
	}

	if err := s.paymentService.CreateTransaction(ctx, transaction); err != nil {
		return 0, fmt.Errorf("failed to create payment transaction: %w", err)
	}

	return id, nil
}

func (s *Service) GetOrder(ctx context.Context, id int) (*entities.Order, error) {
	if id <= 0 {
		return nil, entities.ErrInvalidOrderData
	}

	orderObj, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrOrderNotFound) {
			return nil, entities.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return orderObj, nil
}

func (s *Service) GetOrdersByUserID(ctx context.Context, userID int) ([]entities.Order, error) {
	if userID <= 0 {
		return nil, ErrInvalidOrderData
	}

	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return orders, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return ErrInvalidOrderData
	}

	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	currentOrder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if currentOrder.Status == entities.OrderStatusCompleted ||
		currentOrder.Status == entities.OrderStatusCancelled {
		return ErrOrderAlreadyClosed
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if status == entities.OrderStatusCompleted || status == entities.OrderStatusCancelled {
		if err := s.carService.UpdateStatus(ctx, currentOrder.CarID, "available"); err != nil {
			return fmt.Errorf("failed to update car status: %w", err)
		}
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidOrderData
	}

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == entities.OrderStatusCompleted ||
		order.Status == entities.OrderStatusCancelled {
		return ErrOrderAlreadyClosed
	}

	if err := s.repo.UpdateStatus(ctx, id, entities.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if err := s.carService.UpdateStatus(ctx, order.CarID, "available"); err != nil {
		return fmt.Errorf("failed to update car status: %w", err)
	}

	if err := s.userService.DeductBalance(ctx, order.UserID, -order.TotalPrice); err != nil {
		return fmt.Errorf("failed to refund user balance: %w", err)
	}

	return nil
}

func (s *Service) ListAllOrders(ctx context.Context) ([]entities.Order, error) {
	orders, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

func validateOrder(o *entities.Order) error {
	if o == nil {
		return errors.New("order is nil")
	}
	if o.UserID <= 0 {
		return errors.New("invalid user ID")
	}
	if o.CarID <= 0 {
		return errors.New("invalid car ID")
	}
	if o.TotalPrice <= 0 {
		return errors.New("total price must be positive")
	}
	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case entities.OrderStatusPending,
		entities.OrderStatusConfirmed,
		entities.OrderStatusCompleted,
		entities.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

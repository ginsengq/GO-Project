package entities

import (
	"errors"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CarID      int       `json:"car_id"`
	Status     string    `json:"status"`
	Deposit    float64   `json:"deposit"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
	OrderStatusCompleted = "completed"
	OrderStatusConfirmed = "confirmed"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderData   = errors.New("invalid order data")
	ErrOrderAlreadyClosed = errors.New("order is already completed or cancelled")
	ErrInvalidStatus      = errors.New("invalid order status")
)

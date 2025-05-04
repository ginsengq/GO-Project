package entities

import (
	"errors"
	"time"
)

var (
	ErrInvalidID    = errors.New("invalid ID")
	ErrInvalidInput = errors.New("invalid input")
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Password     string    `json:"password"`
	Balance      float64   `json:"balance"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var ErrNotFoundund = errors.New("user not found")

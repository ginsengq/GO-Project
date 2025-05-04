package entities

import (
	"errors"
	"time"
)

type Car struct {
	ID        int       `json:"id" db:"id"`
	Brand     string    `json:"brand" db:"brand"`
	Model     string    `json:"model" db:"model"`
	Year      int       `json:"year" db:"year"`
	Price     float64   `json:"price" db:"price"`
	Mileage   int       `json:"mileage" db:"mileage"`
	Color     string    `json:"color" db:"color"`
	Status    CarStatus `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CarStatus string

const (
	CarStatusAvailable CarStatus = "available"
	CarStatusReserved  CarStatus = "reserved"
	CarStatusSold      CarStatus = "sold"
)

type CarFilter struct {
	Brand     *string  `json:"brand,omitempty"`
	Model     *string  `json:"model,omitempty"`
	YearFrom  *int     `json:"year_from,omitempty"`
	YearTo    *int     `json:"year_to,omitempty"`
	MinPrice  *float64 `json:"min_price,omitempty"`
	MaxPrice  *float64 `json:"max_price,omitempty"`
	Status    *string  `json:"status,omitempty"`
	Color     *string  `json:"color,omitempty"`
	Limit     *int     `json:"limit,omitempty"`
	Offset    *int     `json:"offset,omitempty"`
	SortBy    *string  `json:"sort_by,omitempty"`
	SortOrder *string  `json:"sort_order,omitempty"`
}

type CarUpdate struct {
	Brand   *string    `json:"brand,omitempty"`
	Model   *string    `json:"model,omitempty"`
	Year    *int       `json:"year,omitempty"`
	Price   *float64   `json:"price,omitempty"`
	Mileage *int       `json:"mileage,omitempty"`
	Color   *string    `json:"color,omitempty"`
	Status  *CarStatus `json:"status,omitempty"`
}

var (
	ErrNotFound = errors.New("car not found")
)

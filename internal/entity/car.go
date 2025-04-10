package entity

import "time"

type CarStatus string

const (
	Available CarStatus = "available"
	Reserved  CarStatus = "reserved"
	Sold      CarStatus = "sold"
)

type Car struct {
	ID        int64     `json:"id"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Year      string    `json:"year"`
	Price     float64   `json:"price"`
	Mileage   float64   `json:"mileage"`
	Color     string    `json:"color"`
	Status    CarStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

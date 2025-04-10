package entity

import "time"

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

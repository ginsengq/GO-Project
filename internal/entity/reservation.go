package entity

import "time"

type Reservation struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CarID     int       `json:"car_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

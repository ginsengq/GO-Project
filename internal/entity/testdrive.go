package entity

import "time"

type TestDrive struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CarID     int       `json:"car_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

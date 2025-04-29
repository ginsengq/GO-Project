package services

import (
	repository "myproject/internal/repositories"
)

type Service struct {
	Auth AuthService
	// другие сервисы, например:
	// Car   CarService
	// Order OrderService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.User),
	}
}

package services

import (
	"errors"
	"myproject/internal/entity"
	repository "myproject/internal/repositories"
)

type AuthService interface {
	Register(name, email, password string) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(name, email, password string) error {
	return s.repo.CreateUser(entity.User{Name: name, Email: email, PasswordHash: password})
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil || user.PasswordHash != password {
		return "", errors.New("invalid credentials")
	}

	return "mock-jwt-token", nil
}

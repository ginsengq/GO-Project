package userservice

import (
	"context"
	"errors"
	"fmt"
	"myproject/internal/entities"
	userrepo "myproject/internal/repositories/user"

	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	Authenticate(ctx context.Context, email, password string) (string, *entities.User, error)
	Count(ctx context.Context) (int, error)
	CheckBalance(ctx context.Context, userID int, amount float64) (bool, error)
	DeductBalance(ctx context.Context, userID int, amount float64) error
}

type Service struct {
	repo userrepo.Repository
}

func NewUserService(repo userrepo.Repository) *Service {
	return &Service{repo: repo}
}

func generateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) Create(ctx context.Context, user *entities.User) error {
	hashedPassword, err := generateHash(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword
	user.Balance = 0
	user.Role = "customer"

	return s.repo.Create(ctx, user)
}

func (s *Service) GetByID(ctx context.Context, id int) (*entities.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) Update(ctx context.Context, user *entities.User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !checkPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := generateHash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.PasswordHash = hashedPassword
	return s.repo.Update(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, *entities.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("incorrect password")
	}

	token := "generated_token"
	return token, user, nil
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s *Service) CheckBalance(ctx context.Context, userID int, amount float64) (bool, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}
	return user.Balance >= amount, nil
}

func (s *Service) DeductBalance(ctx context.Context, userID int, amount float64) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Balance < amount {
		return errors.New("insufficient funds")
	}

	newBalance := user.Balance - amount
	return s.repo.UpdateBalance(ctx, userID, newBalance)
}

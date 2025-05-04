package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"myproject/internal/entities"
	userrepo "myproject/internal/repositories/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	Authenticate(ctx context.Context, email, password string) (string, *entities.User, error) // Возвращаем JWT
	Count(ctx context.Context) (int, error)
	GenerateJWT(user *entities.User) (string, error) // добавлено в UserService
}

type AuthService interface {
	ValidateJWT(tokenString string) (*jwt.MapClaims, error)
}

type userService struct {
	repo      userrepo.Repository
	jwtSecret []byte
}

type authService struct {
	jwtSecret []byte
}

func NewUserService(repo userrepo.Repository) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: []byte("your_jwt_secret"),
	}
}

func NewAuthService(jwtSecret string) AuthService {
	return &authService{jwtSecret: []byte(jwtSecret)}
}

func (s *userService) Create(ctx context.Context, user *entities.User) error {
	if user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("invalid user data")
	}

	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, entities.ErrNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id int) (*entities.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *entities.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid ID")
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if limit <= 0 || offset < 0 {
		return nil, errors.New("invalid pagination parameters")
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *userService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(ctx, user)
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (string, *entities.User, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *userService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s *userService) GenerateJWT(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	} else {
		return nil, errors.New("invalid token claims")
	}
}

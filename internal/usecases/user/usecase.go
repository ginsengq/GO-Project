package usercase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"myproject/internal/entities"

	repoUser "myproject/internal/repositories/user"

	"github.com/golang-jwt/jwt/v5"
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
	Authenticate(ctx context.Context, email, password string) (string, *entities.User, error) // Возвращаем JWT
	Count(ctx context.Context) (int, error)
}

type useCase struct {
	repo      repoUser.Repository
	jwtSecret []byte
}

func New(repo repoUser.Repository, jwtSecret string) UseCase {
	return &useCase{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (u *useCase) Create(ctx context.Context, user *entities.User) error {
	if user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("invalid user data")
	}

	existingUser, err := u.repo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, entities.ErrNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	return u.repo.Create(ctx, user)
}

func (u *useCase) GetByID(ctx context.Context, id int) (*entities.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return u.repo.GetByEmail(ctx, email)
}

func (u *useCase) Update(ctx context.Context, user *entities.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	existingUser, err := u.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return u.repo.Update(ctx, user)
}

func (u *useCase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid ID")
	}
	return u.repo.Delete(ctx, id)
}

func (u *useCase) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if limit <= 0 || offset < 0 {
		return nil, errors.New("invalid pagination parameters")
	}
	return u.repo.List(ctx, limit, offset)
}

func (u *useCase) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}

	user, err := u.repo.GetByID(ctx, userID)
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
	return u.repo.Update(ctx, user)
}

func (u *useCase) Authenticate(ctx context.Context, email, password string) (string, *entities.User, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, user, nil
}

func (u *useCase) Count(ctx context.Context) (int, error) {
	return u.repo.Count(ctx)
}

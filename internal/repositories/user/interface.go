package user

import "myproject/internal/entity"

type UserRepository interface {
	Create(user entity.User) (int, error)
	GetByEmail(email string) (*entity.User, error)
	GetByID(id int) (*entity.User, error)
	Update(user entity.User) error
	Delete(id int) error
}

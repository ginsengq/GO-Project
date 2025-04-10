package service

import "myproject/internal/repository"

type Car interface {
}

type CarUseCase struct {
}

func NewCarUseCase() {

}

type Service struct {
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}

package service

import (
	apiserver "mongo_db"
	"mongo_db/pkg/repository"
)

type Authorization interface {
	CreateUser(user apiserver.User) error
	GenerateRefreshToken() (string, error)
	GenerateAccessToken(id string) (string, error)
	GetUserById(id string) (*apiserver.User, error)
	UpdateRefreshToken(id string, refreshToken string) error
	ParseAccessToken(accessToken string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthSerivce(repo.Authorization),
	}
}

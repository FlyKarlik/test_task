package repository

import (
	apiserver "mongo_db"

	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(user apiserver.User) error
	GetUserById(id string) (*apiserver.User, error)
	UpdateRefreshToken(id string, refreshToken string) error
}

type Repository struct {
	Authorization
}

func NewUserRepostiroy(db *mongo.Collection) *Repository {
	return &Repository{
		Authorization: NewAuthMongoDb(db),
	}
}

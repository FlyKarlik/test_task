package repository

import (
	"context"
	apiserver "mongo_db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongoDb struct {
	db *mongo.Collection
}

func NewAuthMongoDb(db *mongo.Collection) *AuthMongoDb {
	return &AuthMongoDb{db: db}
}

func (a *AuthMongoDb) CreateUser(user apiserver.User) error {
	_, err := a.db.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthMongoDb) UpdateRefreshToken(id string, refreshToken string) error {
	_, err := a.db.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": bson.M{"refresh_token": refreshToken}})
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthMongoDb) GetUserById(id string) (*apiserver.User, error) {
	var user apiserver.User
	if err := a.db.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

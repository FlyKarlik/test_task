package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	User         string
	Password     string
	DatabaseName string
	Collection   string
	Port         string
}

func NewMongoDb(cfg Config) (*mongo.Collection, error) {
	databaseURL := fmt.Sprintf("mongodb://%s:%s@localhost:%s", cfg.User, cfg.Password, cfg.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create a client: %s", err.Error())
	}

	if err = client.Connect(context.TODO()); err != nil {
		return nil, fmt.Errorf("failed to create a database connection: %s", err.Error())
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("the connection to the database could not be verifed: %s", err.Error())
	}

	return client.Database(cfg.DatabaseName).Collection(cfg.Collection), nil

}

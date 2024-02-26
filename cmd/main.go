package main

import (
	"log"
	apiserver "mongo_db"
	"mongo_db/pkg/handler"
	"mongo_db/pkg/repository"
	"mongo_db/pkg/service"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatalf("error initilizing config file: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	mongodb, err := repository.NewMongoDb(repository.Config{
		Password:     os.Getenv("PASSWORD_DB"),
		DatabaseName: viper.GetString("db.database_name"),
		Collection:   viper.GetString("db.collection"),
		User:         viper.GetString("db.user"),
		Port:         viper.GetString("db.port"),
	})
	if err != nil {
		log.Fatalf("failed to initilize database: %s", err.Error())
	}
	repository := repository.NewUserRepostiroy(mongodb)
	services := service.NewService(repository)
	handlers := handler.NewHandler(services)
	server := apiserver.NewServer()

	if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error for running http server: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

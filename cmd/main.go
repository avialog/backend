package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/controller"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Panic("Error loading .env file")
	}

	cfg := config.NewConfig()
	decodedFirebaseKey, err := cfg.DecodeFirebaseKey()
	if err != nil {
		logrus.Panic(err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(decodedFirebaseKey))
	if err != nil {
		logrus.Panic(err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		logrus.Panic(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}

	server := gin.Default()

	repositories, err := repository.NewRepositories(db)
	if err != nil {
		logrus.Panic(err)
	}
	services := service.NewServices(repositories, cfg, utils.GetValidator(), authClient)
	controllers := controller.NewControllers(services, cfg)
	controllers.Route(server)

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	err = server.Run(":" + port)
	if err != nil {
		logrus.Panic(err)
	}
}

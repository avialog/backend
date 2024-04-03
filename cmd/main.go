package main

import (
	"context"
	firebase "firebase.google.com/go"
	_ "github.com/avialog/backend/docs"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/controller"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// @title Avialog API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @host localhost:3000
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Authorization by JWT token

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

	//configure swagger endpoint
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = server.Run(":" + port)
	if err != nil {
		logrus.Panic(err)
	}
}

package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/avialog/backend/internal/controller"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	logrus.Panic("Error loading .env file")
	//}

	config := dto.Config{
		DSN: os.Getenv("DSN"),
	}

	opt := option.WithCredentialsFile("cmd/firebase.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Panic(err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		logrus.Panic(err)
	}

	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=testserver dbname=newDb port=5432 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}

	server := gin.Default()

	repositories, err := repository.NewRepositories(db)
	if err != nil {
		logrus.Panic(err)
	}
	services := service.NewServices(repositories, config, utils.GetValidator())
	controllers := controller.NewControllers(services, config, authClient)
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

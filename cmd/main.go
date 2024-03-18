package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := dto.Config{
		DSN: os.Getenv("DSN"),
	}

	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}

	server := gin.Default()

	repositories, err := repository.NewRepositories(db)
	if err != nil {
		logrus.Panic(err)
	}

	_ = repositories

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	err = server.Run(":" + port)
	if err != nil {
		logrus.Panic(err)
	}
}

package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

func main() {
	config := dto.Config{Port: 3000}

	db, err := gorm.Open(postgres.Open("postgresql://localhost:5432"), &gorm.Config{}) //assign official db in future
	if err != nil {
		logrus.Panic(err)
	}

	server := gin.Default()

	repositories, err := repository.NewRepositories(db)
	if err != nil {
		logrus.Panic(err)
	}

	_ = repositories

	err = server.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		logrus.Panic(err)
	}
}

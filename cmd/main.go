package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

func main() {
	config := dto.Config{Port: 8080}

	_, err := gorm.Open(postgres.Open("postgresql://localhost:5432"), &gorm.Config{}) //assign official db in future

	if err != nil {
		panic(err)
	}

	server := gin.Default()

	//repositories := repository.NewRepositories(db) // To uncomment in the future

	err = server.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}
}

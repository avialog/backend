package main

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	config := dto.Config{}
	server := gin.Default()
	err := server.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}
}

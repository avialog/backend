package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoController interface {
	Info(*gin.Context)
}

type infoController struct{}

func newInfoController() InfoController {
	return &infoController{}
}

func (*infoController) Info(c *gin.Context) {
	c.JSON(http.StatusOK, dto.ServerInfo{Healthy: true})
}

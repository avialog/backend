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

// Info godoc
//
// @Summary Get server health information
// @Description Get information about the server's health status
// @Tags info
// @Accept  json
// @Produce  json
// @Success 200 {object}      dto.ServerInfo
// @Router  /info [get]
func (*infoController) Info(c *gin.Context) {
	c.JSON(http.StatusOK, dto.ServerInfo{Healthy: true})
}

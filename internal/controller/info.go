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
// @Summary Health check endpoint
// @Description Returns the health status of the server
// @Tags info
// @Produce json
// @Success 200 {object} dto.ServerInfo
// @Router /healthz [get]
func (*infoController) Info(c *gin.Context) {
	c.JSON(http.StatusOK, dto.ServerInfo{Healthy: true})
}

package controller

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AircraftController interface {
	GetAircraft(*gin.Context)
	InsertAircraft(*gin.Context)
	UpdateAircraft(*gin.Context)
	DeleteAircraft(*gin.Context)
}

type aircraftController struct {
	aircraftService service.AircraftService
}

func newAircraftController(aircraftService service.AircraftService) AircraftController {
	return &aircraftController{aircraftService: aircraftService}
}

func (a *aircraftController) GetAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	aircraft, err := a.aircraftService.GetUserAircraft(userID)

	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			util.NewError(ctx, http.StatusBadRequest, err) //TODO: ask about this
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, aircraft)
}

func (a *aircraftController) InsertAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")
}

func (a *aircraftController) UpdateAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")
}

func (a *aircraftController) DeleteAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")
}

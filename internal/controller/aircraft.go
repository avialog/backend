package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
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

// GetAircraft godoc
// @Summary Get user aircraft (all)
// @Description Get user aircraft
// @Tags aircraft
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.AircraftResponse
// @Router /aircraft [get]
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
func (a *aircraftController) GetAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	aircraft, err := a.aircraftService.GetUserAircraft(userID)

	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	aircraftResponse := a.adaptAircraftSlice(aircraft)

	ctx.JSON(http.StatusOK, aircraftResponse)
}

// InsertAircraft godoc
// @Summary Insert aircraft
// @Description Insert aircraft
// @Tags aircraft
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} dto.AircraftResponse
// @Router /aircraft [post]
// @Param aircraft body dto.AircraftRequest true "Aircraft"
// @Failure 400 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
func (a *aircraftController) InsertAircraft(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var aircraftRequest dto.AircraftRequest
	if err := ctx.ShouldBindJSON(&aircraftRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	aircraft, err := a.aircraftService.InsertAircraft(userID, aircraftRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	aircraftResponse := a.adaptAircraft(aircraft)

	ctx.JSON(http.StatusCreated, aircraftResponse)
}

// UpdateAircraft godoc
// @Summary Update aircraft
// @Description Update aircraft
// @Tags aircraft
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.AircraftResponse
// @Router /aircraft [put]
// @Param id path string true "Aircraft ID"
// @Param aircraft body dto.AircraftRequest true "Aircraft"
// @Failure 400 {object} util.HTTPError
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
func (a *aircraftController) UpdateAircraft(ctx *gin.Context) {
	aircraftID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString("userID")

	var aircraftRequest dto.AircraftRequest
	if err := ctx.ShouldBindJSON(&aircraftRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	aircraft, err := a.aircraftService.UpdateAircraft(userID, uint(aircraftID), aircraftRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	aircraftResponse := a.adaptAircraft(aircraft)

	ctx.JSON(http.StatusOK, aircraftResponse)
}

// DeleteAircraft godoc
// @Summary Delete aircraft
// @Description Delete aircraft
// @Tags aircraft
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{message=string} "Aircraft deleted successfully"
// @Router /aircraft [delete]
// @Param id path string true "Aircraft ID"
// @Failure 400 {object} util.HTTPError
// @Failure 404 {object} util.HTTPError
// @Failure 500 {object} util.HTTPError
// @Failure 409 {object} util.HTTPError
func (a *aircraftController) DeleteAircraft(ctx *gin.Context) {
	aircraftID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	userID := ctx.GetString("userID")

	err = a.aircraftService.DeleteAircraft(userID, uint(aircraftID))
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		} else if errors.Is(err, dto.ErrConflict) {
			util.NewError(ctx, http.StatusConflict, err)
			return
		} else if errors.Is(err, dto.ErrNotFound) {
			util.NewError(ctx, http.StatusNotFound, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Aircraft deleted successfully"})
}

func (a *aircraftController) adaptAircraft(aircraft model.Aircraft) dto.AircraftResponse {
	return dto.AircraftResponse{
		ID:                 aircraft.ID,
		AircraftModel:      aircraft.AircraftModel,
		RegistrationNumber: aircraft.RegistrationNumber,
		IsSingleEngine:     aircraft.IsSingleEngine,
		ImageURL:           aircraft.ImageURL,
		Remarks:            aircraft.Remarks,
	}
}

func (a *aircraftController) adaptAircraftSlice(aircraftSlice []model.Aircraft) []dto.AircraftResponse {
	aircraftResponses := make([]dto.AircraftResponse, 0)
	for _, aircraft := range aircraftSlice {
		aircraftResponses = append(aircraftResponses, a.adaptAircraft(aircraft))
	}
	return aircraftResponses
}

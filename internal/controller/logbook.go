package controller

import (
	"errors"
	"github.com/avialog/backend/internal/common"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LogbookController interface {
	GetLogbookEntries(*gin.Context)
	InsertLogbookEntry(*gin.Context)
	UpdateLogbookEntry(*gin.Context)
	DeleteLogbookEntry(*gin.Context)
}

type logbookController struct {
	logbookService service.LogbookService
}

func newLogbookController(logbookService service.LogbookService) LogbookController {
	return &logbookController{logbookService: logbookService}
}

// GetLogbookEntries godoc
//
// @Summary Get user logbook entries
// @Description Get a list of logbook entries for a user
// @Tags logbook
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} []dto.LogbookResponse
// @Failure 500 {object}      util.HTTPError
// @Failure 400 {object}      util.HTTPError
// @Router  /logbook [get]
func (c *logbookController) GetLogbookEntries(ctx *gin.Context) {
	userID := ctx.GetString(common.UserID)
	var start time.Time
	var end time.Time

	var getLogbookRequest dto.GetLogbookRequest
	if err := ctx.ShouldBindJSON(&getLogbookRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if getLogbookRequest.Start == nil && getLogbookRequest.End != nil || getLogbookRequest.Start != nil && getLogbookRequest.End == nil {
		util.NewError(ctx, http.StatusBadRequest, errors.New("both start and end time must be provided or neither"))
		return
	} else if getLogbookRequest.Start == nil && getLogbookRequest.End == nil {
		start = time.Now().AddDate(0, 0, -90)
		end = time.Now()
	} else {
		start = time.Unix(*getLogbookRequest.Start, 0)
		end = time.Unix(*getLogbookRequest.End, 0)
	}
	flights, err := c.logbookService.GetLogbookEntries(userID, start, end)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, flights)
}

// InsertLogbookEntry godoc
//
// @Summary Insert a new logbook entry
// @Description Insert a new logbook entry for a user
// @Tags logbook
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   logbookRequest     body     dto.LogbookRequest true    "Logbook entry information to insert"
// @Success 201 {object}      dto.LogbookResponse
// @Failure 400 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /logbook [post]
func (c *logbookController) InsertLogbookEntry(ctx *gin.Context) {
	userID := ctx.GetString(common.UserID)

	var logbookRequest dto.LogbookRequest
	if err := ctx.ShouldBindJSON(&logbookRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	logbookResponse, err := c.logbookService.InsertLogbookEntry(userID, logbookRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, logbookResponse)
}

// UpdateLogbookEntry godoc
//
// @Summary Update an existing logbook entry
// @Description Update an existing logbook entry for a user
// @Tags logbook
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id                path     int        true        "Flight ID to update"
// @Param   logbookRequest     body     dto.LogbookRequest true    "Logbook entry information to update"
// @Success 200 {object}      dto.LogbookResponse
// @Failure 400 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router /logbook/{id} [put]
func (c *logbookController) UpdateLogbookEntry(ctx *gin.Context) {
	flightID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString(common.UserID)

	var logbookRequest dto.LogbookRequest
	if err := ctx.ShouldBindJSON(&logbookRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	logbookResponse, err := c.logbookService.UpdateLogbookEntry(userID, uint(flightID), logbookRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, logbookResponse)
}

// DeleteLogbookEntry godoc
//
// @Summary Delete an existing logbook entry
// @Description Delete an existing logbook entry for a user
// @Tags logbook
// @Produce  json
// @Security ApiKeyAuth
// @Param   id                path     int        true        "Flight ID to delete"
// @Success 200 {object}      object{message=string} "Logbook entry deleted successfully"
// @Failure 400 {object}      util.HTTPError
// @Failure 404 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /logbook/{id} [delete]
func (c *logbookController) DeleteLogbookEntry(ctx *gin.Context) {
	flightID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString(common.UserID)

	err = c.logbookService.DeleteLogbookEntry(userID, uint(flightID))
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			util.NewError(ctx, http.StatusNotFound, err)
			return
		} else if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logbook entry deleted successfully"})
}

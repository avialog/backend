package controller

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @Summary Get user flights
// @Description Get a list of flights for a user
// @Tags flights
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   userRequest       body     dto.GetFlightsRequest true       "start and end time to get flights for user"
// @Success 200 {object}      dto.LogbookResponse
// @Failure 500 {object}      util.HTTPError
// @Router  /flights [get]
func (c *logbookController) GetLogbookEntries(ctx *gin.Context) {
	userID := ctx.GetString("userID")
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
		start = time.Now()
		end = start.AddDate(0, -3, 0)
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

func (c *logbookController) InsertLogbookEntry(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var logbookRequest dto.LogbookRequest
	if err := ctx.ShouldBindJSON(&logbookRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	logbookResponse, err := c.logbookService.InsertLogbookEntry(userID, logbookRequest)
	if err != nil {

	}
}

func (c *logbookController) UpdateLogbookEntry(ctx *gin.Context) {

}

func (c *logbookController) DeleteLogbookEntry(ctx *gin.Context) {

}

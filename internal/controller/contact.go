package controller

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ContactController interface {
	GetContacts(*gin.Context)
	InsertContact(*gin.Context)
	UpdateContact(*gin.Context)
	DeleteContact(*gin.Context)
}

type contactController struct {
	contactService service.ContactService
}

func newContactController(contactService service.ContactService) ContactController {
	return &contactController{contactService: contactService}
}

func (c *contactController) GetContacts(ctx *gin.Context) {
	// TODO: add getting user id from JWT token
	userID := uint(1)

	contacts, err := c.contactService.GetUserContacts(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

func (c *contactController) InsertContact(ctx *gin.Context) {
	// TODO: add getting user id from JWT token
	userID := uint(1)

	var contactRequest dto.ContactRequest
}

func (c *contactController) UpdateContact(ctx *gin.Context) {
	// TODO: add getting user id from JWT token
	userID := uint(1)

	var contactRequest dto.ContactRequest
}

func (c *contactController) DeleteContact(ctx *gin.Context) {
	//parse to uint32
	contactID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: add getting user id from JWT token
	userID := uint(1)

	err = c.contactService.DeleteContact(userID, uint(contactID))
}

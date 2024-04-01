package controller

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
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
	userID := ctx.GetString("userID")

	contacts, err := c.contactService.GetUserContacts(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contactsResponse := c.adaptContacts(contacts)

	ctx.JSON(http.StatusOK, contactsResponse)
}

func (c *contactController) InsertContact(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var contactRequest dto.ContactRequest
	if err := ctx.ShouldBindJSON(&contactRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := c.contactService.InsertContact(userID, contactRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contactResponse := c.adaptContact(contact)

	ctx.JSON(http.StatusCreated, contactResponse)
}

func (c *contactController) UpdateContact(ctx *gin.Context) {
	contactID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetString("userID")

	var contactRequest dto.ContactRequest
	if err := ctx.ShouldBindJSON(&contactRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := c.contactService.UpdateContact(userID, uint(contactID), contactRequest)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, dto.ErrBadRequest) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contactResponse := c.adaptContact(contact)

	ctx.JSON(http.StatusOK, contactResponse)
}

func (c *contactController) DeleteContact(ctx *gin.Context) {
	contactID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetString("userID")

	err = c.contactService.DeleteContact(userID, uint(contactID))
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

func (c *contactController) adaptContact(contact model.Contact) dto.ContactResponse {
	return dto.ContactResponse{
		AvatarURL:    contact.AvatarURL,
		FirstName:    contact.FirstName,
		LastName:     contact.LastName,
		Company:      contact.Company,
		Phone:        contact.Phone,
		EmailAddress: contact.EmailAddress,
		Note:         contact.Note,
	}
}

func (c *contactController) adaptContacts(contacts []model.Contact) []dto.ContactResponse {
	var contactsResponse []dto.ContactResponse
	for _, contact := range contacts {
		contactsResponse = append(contactsResponse, c.adaptContact(contact))
	}
	return contactsResponse
}

package controller

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
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

// GetContacts godoc
//
// @Summary Get user contacts
// @Description Get a list of contacts for a user
// @Tags contacts
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array}       dto.ContactResponse
// @Failure 500 {object}      util.HTTPError
// @Router  /contacts [get]
func (c *contactController) GetContacts(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	contacts, err := c.contactService.GetUserContacts(userID)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	contactsResponse := c.adaptContacts(contacts)

	ctx.JSON(http.StatusOK, contactsResponse)
}

// InsertContact godoc
//
// @Summary Insert a new contact
// @Description Insert a new contact for a user
// @Tags contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   contactRequest    body     dto.ContactRequest true    "Contact information to insert"
// @Success 201 {object}      dto.ContactResponse
// @Failure 400 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /contacts [post]
func (c *contactController) InsertContact(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var contactRequest dto.ContactRequest
	if err := ctx.ShouldBindJSON(&contactRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	contact, err := c.contactService.InsertContact(userID, contactRequest)
	if err != nil {
		if errors.Is(err, dto.ErrBadRequest) {
			util.NewError(ctx, http.StatusBadRequest, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	contactResponse := c.adaptContact(contact)

	ctx.JSON(http.StatusCreated, contactResponse)
}

// UpdateContact godoc
//
// @Summary Update an existing contact
// @Description Update an existing contact for a user
// @Tags contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id                path     int        true        "Contact ID to update"
// @Param   contactRequest    body     dto.ContactRequest true    "Contact information to update"
// @Success 200 {object}      dto.ContactResponse
// @Failure 400 {object}      util.HTTPError
// @Failure 404 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /contacts/{id} [put]
func (c *contactController) UpdateContact(ctx *gin.Context) {
	contactID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString("userID")

	var contactRequest dto.ContactRequest
	if err := ctx.ShouldBindJSON(&contactRequest); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	contact, err := c.contactService.UpdateContact(userID, uint(contactID), contactRequest)
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

	contactResponse := c.adaptContact(contact)

	ctx.JSON(http.StatusOK, contactResponse)
}

// DeleteContact godoc
//
// @Summary Delete an existing contact
// @Description Delete an existing contact for a user
// @Tags contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id                path     int        true        "Contact ID to delete"
// @Success 200 {object}      object{message=string} "Contact deleted successfully"
// @Failure 400 {object}      util.HTTPError
// @Failure 404 {object}      util.HTTPError
// @Failure 500 {object}      util.HTTPError
// @Router  /contacts/{id} [delete]
func (c *contactController) DeleteContact(ctx *gin.Context) {
	contactID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString("userID")

	err = c.contactService.DeleteContact(userID, uint(contactID))
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			util.NewError(ctx, http.StatusNotFound, err)
			return
		}
		util.NewError(ctx, http.StatusInternalServerError, err)
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

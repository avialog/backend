package service

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -source=contact.go -destination=contact_mock.go -package service
type ContactService interface {
	InsertContact(userID string, contactRequest dto.ContactRequest) (model.Contact, error)
	GetUserContacts(userID string) ([]model.Contact, error)
	UpdateContact(userID string, id uint, contactRequest dto.ContactRequest) (model.Contact, error)
	DeleteContact(userID string, id uint) error
}

type contactService struct {
	contactRepository repository.ContactRepository
	validator         *validator.Validate
	config            config.Config
}

func newContactService(contactRepository repository.ContactRepository, config config.Config, validator *validator.Validate) ContactService {
	return &contactService{contactRepository, validator, config}
}

func (c *contactService) InsertContact(userID string, contactRequest dto.ContactRequest) (model.Contact, error) {
	contact := model.Contact{
		UserID:       userID,
		FirstName:    contactRequest.FirstName,
		LastName:     contactRequest.LastName,
		Phone:        contactRequest.Phone,
		AvatarURL:    contactRequest.AvatarURL,
		Company:      contactRequest.Company,
		EmailAddress: contactRequest.EmailAddress,
		Note:         contactRequest.Note,
	}

	err := c.validator.Struct(contact)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return model.Contact{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return model.Contact{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}
	return c.contactRepository.Create(contact)
}

func (c *contactService) GetUserContacts(userID string) ([]model.Contact, error) {
	return c.contactRepository.GetByUserID(userID)
}

func (c *contactService) DeleteContact(userID string, id uint) error {
	err := c.contactRepository.DeleteByUserIDAndID(userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *contactService) UpdateContact(userID string, id uint, contactRequest dto.ContactRequest) (model.Contact, error) {
	contact, err := c.contactRepository.GetByUserIDAndID(userID, id)
	if err != nil {
		return model.Contact{}, err
	}

	contact.FirstName = contactRequest.FirstName
	contact.LastName = contactRequest.LastName
	contact.Phone = contactRequest.Phone
	contact.AvatarURL = contactRequest.AvatarURL
	contact.Company = contactRequest.Company
	contact.EmailAddress = contactRequest.EmailAddress
	contact.Note = contactRequest.Note

	err = c.validator.Struct(contact)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return model.Contact{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return model.Contact{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}

	return c.contactRepository.Save(contact)
}

package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
)

type ContactService interface {
	InsertContact(userID uint, contactRequest dto.ContactRequest) (model.Contact, error)
	GetUserContacts(userID uint) ([]model.Contact, error)
	UpdateContact(userID, id uint, contactRequest dto.ContactRequest) (model.Contact, error)
	DeleteContact(userID, id uint) error
}

type contactService struct {
	contactRepository repository.ContactRepository
	config            dto.Config
}

func newContactService(contactRepository repository.ContactRepository, config dto.Config) ContactService {
	return &contactService{contactRepository, config}
}

func (c contactService) InsertContact(userID uint, contactRequest dto.ContactRequest) (model.Contact, error) {
	var contact model.Contact
	contact.UserID = userID
	contact.FirstName = contactRequest.FirstName
	contact.LastName = contactRequest.LastName
	contact.Phone = contactRequest.Phone
	contact.AvatarURL = contactRequest.AvatarURL
	contact.Company = contactRequest.Company
	contact.EmailAddress = contactRequest.EmailAddress
	contact.Note = contactRequest.Note
	return c.contactRepository.Save(contact)
}

func (c contactService) GetUserContacts(userID uint) ([]model.Contact, error) {
	return c.contactRepository.GetByUserID(userID)
}

func (c contactService) DeleteContact(userID, id uint) error {
	rowsAffected, err := c.contactRepository.DeleteByUserIDAndID(userID, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no contact to delete or unauthorized to delete contact")
	}
	return nil
}

func (c contactService) UpdateContact(userID, id uint, contactRequest dto.ContactRequest) (model.Contact, error) {
	contact, err := c.contactRepository.GetByUserIDAndID(userID, id)
	if err != nil {
		return model.Contact{}, err
	}

	if contact.UserID != userID {
		return model.Contact{}, errors.New("unauthorized to update contact")
	}

	contact.FirstName = contactRequest.FirstName
	contact.LastName = contactRequest.LastName
	contact.Phone = contactRequest.Phone
	contact.AvatarURL = contactRequest.AvatarURL
	contact.Company = contactRequest.Company
	contact.EmailAddress = contactRequest.EmailAddress
	contact.Note = contactRequest.Note

	return c.contactRepository.Update(contact)
}
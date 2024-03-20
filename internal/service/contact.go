package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
)

type ContactService interface {
	Insert(userID uint, contactRequest dto.ContactRequest) (model.Contact, error)
	GetByUserID(userID uint) ([]model.Contact, error)
	Update(userID, id uint, contactRequest dto.ContactRequest) (model.Contact, error)
	DeleteByID(userID, id uint) error
}

type contactService struct {
	contactRepository repository.ContactRepository
	config            dto.Config
}

func newContactService(contactRepository repository.ContactRepository, config dto.Config) ContactService {
	return &contactService{contactRepository, config}
}

func (c contactService) Insert(userID uint, contactRequest dto.ContactRequest) (model.Contact, error) {
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

func (c contactService) GetByUserID(userID uint) ([]model.Contact, error) {
	return c.contactRepository.GetByUserID(userID)
}

// Zwraca error jeżeli nie udało się nic usunąć (bo np. nie było co usunąć lub "chcemy usunąć nie nasz kontakt")
func (c contactService) DeleteByID(userID, id uint) error {
	return c.contactRepository.DeleteByID(userID, id)
}

// tu nie da się uniknąć wysłania 2 zapytań
func (c contactService) Update(userID, id uint, contactRequest dto.ContactRequest) (model.Contact, error) {
	contact, err := c.contactRepository.GetByID(id) //pobieram być może obcy rekord
	if err != nil {
		return model.Contact{}, err
	}
	//ale tu sprawdzam czy nie jest obcy, jak jest obcy to wywalam
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

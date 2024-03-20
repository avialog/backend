package repository

import (
	"errors"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type ContactRepository interface {
	Save(contact model.Contact) (model.Contact, error)
	GetByID(id uint) (model.Contact, error)
	GetByUserID(id uint) ([]model.Contact, error)
	Update(contact model.Contact) (model.Contact, error)
	DeleteByID(userID, id uint) error
}

type contact struct {
	db *gorm.DB
}

func newContactRepository(db *gorm.DB) ContactRepository {
	return &contact{
		db: db,
	}
}

func (c contact) Save(contact model.Contact) (model.Contact, error) {
	result := c.db.Create(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c contact) GetByID(id uint) (model.Contact, error) {
	var contact model.Contact
	result := c.db.First(&contact, id)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}
	return contact, nil
}

func (c contact) Update(contact model.Contact) (model.Contact, error) {
	if _, err := c.GetByID(contact.ID); err != nil {
		return model.Contact{}, err
	}

	result := c.db.Save(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

// Wykonujemy usunięcie nawet jak nie ma rekordu
func (c contact) DeleteByID(userID, id uint) error {

	result := c.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Contact{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no contact found or unauthorized")
	}

	return nil
}

func (c contact) GetByUserID(userID uint) ([]model.Contact, error) {
	var contact []model.Contact
	result := c.db.Where("user_id = ?", userID).Find(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return contact, nil
}

package repository

import (
	"errors"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=contact.go -destination=contact_mock.go -package repository
type ContactRepository interface {
	Create(contact model.Contact) (model.Contact, error)
	GetByUserIDAndID(userID, id uint) (model.Contact, error)
	GetByUserID(userID uint) ([]model.Contact, error)
	Save(contact model.Contact) (model.Contact, error)
	DeleteByUserIDAndID(userID, id uint) error
}

type contact struct {
	db *gorm.DB
}

func newContactRepository(db *gorm.DB) ContactRepository {
	return &contact{
		db: db,
	}
}

func (c *contact) Create(contact model.Contact) (model.Contact, error) {
	result := c.db.Create(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c *contact) GetByUserIDAndID(userID, id uint) (model.Contact, error) {
	var contact model.Contact
	result := c.db.Where("user_id = ? AND id = ?", userID, id).First(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c *contact) Save(contact model.Contact) (model.Contact, error) {
	result := c.db.Save(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c *contact) DeleteByUserIDAndID(userID, id uint) error {

	result := c.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("contact cannot be deleted")
	}

	return nil
}

func (c *contact) GetByUserID(userID uint) ([]model.Contact, error) {
	var contact []model.Contact
	result := c.db.Where("user_id = ?", userID).Find(&contact)
	if result.Error != nil {
		return nil, result.Error
	}

	return contact, nil
}

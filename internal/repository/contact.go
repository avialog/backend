package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=contact.go -destination=contact_mock.go -package repository
type ContactRepository interface {
	Save(contact model.Contact) (model.Contact, error)
	GetByUserIDAndID(userID, id uint) (model.Contact, error)
	GetByUserID(userID uint) ([]model.Contact, error)
	Update(contact model.Contact) (model.Contact, error)
	DeleteByUserIDAndID(userID, id uint) (int64, error)
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

func (c contact) GetByUserIDAndID(userID, id uint) (model.Contact, error) {
	var contact model.Contact
	result := c.db.Where("user_id = ? AND id = ?", userID, id).First(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c contact) Update(contact model.Contact) (model.Contact, error) {
	if _, err := c.GetByUserIDAndID(contact.UserID, contact.ID); err != nil {
		return model.Contact{}, err
	}

	result := c.db.Save(&contact)
	if result.Error != nil {
		return model.Contact{}, result.Error
	}

	return contact, nil
}

func (c contact) DeleteByUserIDAndID(userID, id uint) (int64, error) {

	result := c.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Contact{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (c contact) GetByUserID(userID uint) ([]model.Contact, error) {
	var contact []model.Contact
	result := c.db.Where("user_id = ?", userID).Find(&contact)
	if result.Error != nil {
		return nil, result.Error
	}

	return contact, nil
}

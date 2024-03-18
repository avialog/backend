package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user model.User) (model.User, error)
	GetByID(id uint) (model.User, error)
	Update(user model.User) (model.User, error)
	DeleteByID(id uint) error
}

type user struct {
	db *gorm.DB
}

func newUserRepository(db *gorm.DB) UserRepository {
	return &user{
		db: db,
	}
}

func (u user) Save(user model.User) (model.User, error) {
	result := u.db.Create(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (u user) GetByID(id uint) (model.User, error) {
	var user model.User
	result := u.db.First(&user, id)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (u user) Update(user model.User) (model.User, error) {
	if _, err := u.GetByID(user.ID); err != nil {
		return model.User{}, err
	}

	result := u.db.Save(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (u user) DeleteByID(id uint) error {
	if _, err := u.GetByID(id); err != nil {
		return err
	}

	result := u.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

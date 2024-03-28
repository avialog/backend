package repository

import (
	"errors"
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=user.go -destination=user_mock.go -package repository
type UserRepository interface {
	Create(user model.User) (model.User, error)
	GetByID(id uint) (model.User, error)
	Save(user model.User) (model.User, error)
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

func (u *user) Create(user model.User) (model.User, error) {
	result := u.db.Create(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (u *user) GetByID(id uint) (model.User, error) {
	var user model.User
	result := u.db.First(&user, id)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (u *user) Save(user model.User) (model.User, error) {
	result := u.db.Save(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (u *user) DeleteByID(id uint) error {
	result := u.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user cannot be deleted")
	}
	return nil
}

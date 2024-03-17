package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type Repositories interface {
}

type repositories struct {
}

func NewRepositories(db *gorm.DB) Repositories {
	err := db.AutoMigrate(&model.User{}, &model.Aircraft{}, &model.Contact{},
		&model.Flight{}, &model.Landing{}, &model.Passenger{})

	if err != nil {
		panic(err)
	}

	return &repositories{}
}

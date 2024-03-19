package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type Repositories interface {
	User() UserRepository
	Aircraft() AircraftRepository
}

type repositories struct {
	userRepository     UserRepository
	aircraftRepository AircraftRepository
	contactRepository  ContactRepository
}

func NewRepositories(db *gorm.DB) (Repositories, error) {
	err := db.AutoMigrate(&model.User{}, &model.Aircraft{}, &model.Contact{},
		&model.Flight{}, &model.Landing{}, &model.Passenger{})

	if err != nil {
		return nil, err
	}

	return &repositories{
		userRepository:     newUserRepository(db),
		aircraftRepository: newAircraftRepository(db),
		contactRepository:  newContactRepository(db),
	}, nil
}

func (r repositories) User() UserRepository {
	return r.userRepository
}

func (r repositories) Aircraft() AircraftRepository { return r.aircraftRepository }

func (r repositories) Contact() ContactRepository { return r.contactRepository }

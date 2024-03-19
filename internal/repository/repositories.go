package repository

import (
	"github.com/avialog/backend/internal/model"
	"gorm.io/gorm"
)

type Repositories interface {
	User() UserRepository
	Flight() FlightRepository
	Landing() LandingRepository
	Passenger() PassengerRepository
	Aircraft() AircraftRepository
	Contact() ContactRepository
}

type repositories struct {
	userRepository      UserRepository
	passengerRepository PassengerRepository
	aircraftRepository  AircraftRepository
	flightRepository    FlightRepository
	landingRepository   LandingRepository
	contactRepository   ContactRepository
}

func NewRepositories(db *gorm.DB) (Repositories, error) {
	err := db.AutoMigrate(&model.User{}, &model.Aircraft{}, &model.Contact{},
		&model.Flight{}, &model.Landing{}, &model.Passenger{})

	if err != nil {
		return nil, err
	}

	return &repositories{
		userRepository:      newUserRepository(db),
		aircraftRepository:  newAircraftRepository(db),
		flightRepository:    newFlightRepository(db),
		landingRepository:   newLandingRepository(db),
		passengerRepository: newPassengerRepository(db),
		contactRepository:   newContactRepository(db),
	}, nil
}

func (r repositories) User() UserRepository {
	return r.userRepository
}

func (r repositories) Aircraft() AircraftRepository { return r.aircraftRepository }

func (r repositories) Flight() FlightRepository { return r.flightRepository }

func (r repositories) Passenger() PassengerRepository { return r.passengerRepository }

func (r repositories) Contact() ContactRepository { return r.contactRepository }

func (r repositories) Landing() LandingRepository { return r.landingRepository }

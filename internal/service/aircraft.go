package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
)

type AircraftService interface {
	Insert(userID uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	GetByUserID(userID uint) ([]model.Aircraft, error)
	Update(userID, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	DeleteByID(userID, id uint) error
	CountFlightsByID(id uint) (int, error)
}

type aircraftService struct {
	aircraftRepository repository.AircraftRepository
	flightRepository   repository.FlightRepository
	config             dto.Config
}

func newAircraftService(aircraftRepository repository.AircraftRepository, flightRepository repository.FlightRepository,
	config dto.Config) AircraftService {
	return &aircraftService{aircraftRepository: aircraftRepository, flightRepository: flightRepository, config: config}
}

func (a aircraftService) Insert(userID uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
	var aircraft model.Aircraft
	aircraft.UserID = userID
	aircraft.AircraftModel = aircraftRequest.AircraftModel
	aircraft.RegistrationNumber = aircraftRequest.RegistrationNumber
	aircraft.ImageURL = aircraftRequest.ImageURL
	aircraft.Remarks = aircraftRequest.Remarks

	return a.aircraftRepository.Save(aircraft)
}

func (a aircraftService) GetByUserID(userID uint) ([]model.Aircraft, error) {
	return a.aircraftRepository.GetByUserID(userID)
}

func (a aircraftService) Update(userID, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
	aircraft, err := a.aircraftRepository.GetByID(id)
	if err != nil {
		return model.Aircraft{}, err
	}

	if aircraft.UserID != userID {
		return model.Aircraft{}, errors.New("unauthorized to update aircraft")
	}

	aircraft.AircraftModel = aircraftRequest.AircraftModel
	aircraft.RegistrationNumber = aircraftRequest.RegistrationNumber
	aircraft.ImageURL = aircraftRequest.ImageURL
	aircraft.Remarks = aircraftRequest.Remarks

	return a.aircraftRepository.Update(aircraft)
}

func (a aircraftService) DeleteByID(userID, id uint) error {
	flights, err := a.flightRepository.GetByAircraftID(id)
	if len(flights) > 0 {
		return errors.New("the plane has assigned flights")
	}

	aircraft, err := a.aircraftRepository.GetByID(id)
	if err != nil {
		return err
	}

	if aircraft.UserID != userID {
		return errors.New("unauthorized to delete aircraft")
	}

	return a.aircraftRepository.DeleteByID(id)
}

func (a aircraftService) CountFlightsByID(id uint) (int, error) {
	flights, err := a.flightRepository.GetByAircraftID(id)
	if err != nil {
		return 0, err
	}

	numberOfFlights := len(flights)

	return numberOfFlights, nil
}

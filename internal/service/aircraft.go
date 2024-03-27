package service

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

type AircraftService interface {
	InsertAircraft(userID uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	GetUserAircraft(userID uint) ([]model.Aircraft, error)
	UpdateAircraft(userID, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	DeleteAircraft(userID, id uint) error
	CountAircraftFlights(userID, id uint) (int64, error)
}

type aircraftService struct {
	aircraftRepository repository.AircraftRepository
	flightRepository   repository.FlightRepository
	validator          *validator.Validate
	config             dto.Config
}

func newAircraftService(aircraftRepository repository.AircraftRepository, flightRepository repository.FlightRepository,
	config dto.Config, validator *validator.Validate) AircraftService {
	return &aircraftService{aircraftRepository: aircraftRepository, flightRepository: flightRepository, config: config, validator: validator}
}

func (a *aircraftService) InsertAircraft(userID uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
	aircraft := model.Aircraft{
		UserID:             userID,
		AircraftModel:      aircraftRequest.AircraftModel,
		RegistrationNumber: aircraftRequest.RegistrationNumber,
		ImageURL:           aircraftRequest.ImageURL,
		Remarks:            aircraftRequest.Remarks,
	}

	err := a.validator.Struct(aircraft)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return model.Aircraft{}, err
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, vErr := range validationErrors {
				return model.Aircraft{}, fmt.Errorf("invalid data in field: %s", vErr.Field())
			}
		}
	}

	return a.aircraftRepository.Create(aircraft)
}

func (a *aircraftService) GetUserAircraft(userID uint) ([]model.Aircraft, error) {
	return a.aircraftRepository.GetByUserID(userID)
}

func (a *aircraftService) UpdateAircraft(userID, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
	aircraft, err := a.aircraftRepository.GetByUserIDAndID(userID, id)
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

	err = a.validator.Struct(aircraft)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return model.Aircraft{}, err
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, vErr := range validationErrors {
				return model.Aircraft{}, fmt.Errorf("invalid data in field: %s", vErr.Field())
			}
		}
	}

	return a.aircraftRepository.Save(aircraft)
}

func (a *aircraftService) DeleteAircraft(userID, id uint) error {
	numberOfFlights, err := a.flightRepository.CountByAircraftID(userID, id)
	if err != nil {
		return err
	}

	if numberOfFlights > 0 {
		return errors.New("the plane has assigned flights")
	}

	err = a.aircraftRepository.DeleteByUserIDAndID(userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *aircraftService) CountAircraftFlights(userID, id uint) (int64, error) {
	numberOfFlights, err := a.flightRepository.CountByAircraftID(userID, id)
	if err != nil {
		return 0, err
	}

	return numberOfFlights, nil
}

package service

import (
	"errors"
	"fmt"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -source=aircraft.go -destination=aircraft_mock.go -package service
type AircraftService interface {
	InsertAircraft(userID string, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	GetUserAircraft(userID string) ([]model.Aircraft, error)
	UpdateAircraft(userID string, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error)
	DeleteAircraft(userID string, id uint) error
}

type aircraftService struct {
	aircraftRepository repository.AircraftRepository
	flightRepository   repository.FlightRepository
	validator          *validator.Validate
	config             config.Config
}

func newAircraftService(aircraftRepository repository.AircraftRepository, flightRepository repository.FlightRepository,
	config config.Config, validator *validator.Validate) AircraftService {
	return &aircraftService{aircraftRepository: aircraftRepository, flightRepository: flightRepository, config: config, validator: validator}
}

func (a *aircraftService) InsertAircraft(userID string, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
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
			return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return model.Aircraft{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}

	return a.aircraftRepository.Create(aircraft)
}

func (a *aircraftService) GetUserAircraft(userID string) ([]model.Aircraft, error) {
	return a.aircraftRepository.GetByUserID(userID)
}

func (a *aircraftService) UpdateAircraft(userID string, id uint, aircraftRequest dto.AircraftRequest) (model.Aircraft, error) {
	aircraft, err := a.aircraftRepository.GetByUserIDAndID(userID, id)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrBadRequest, err)
		}
		return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}

	aircraft.AircraftModel = aircraftRequest.AircraftModel
	aircraft.RegistrationNumber = aircraftRequest.RegistrationNumber
	aircraft.ImageURL = aircraftRequest.ImageURL
	aircraft.Remarks = aircraftRequest.Remarks

	err = a.validator.Struct(aircraft)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return model.Aircraft{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return model.Aircraft{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}

	return a.aircraftRepository.Save(aircraft)
}

func (a *aircraftService) DeleteAircraft(userID string, id uint) error {
	numberOfFlights, err := a.flightRepository.CountByUserIDAndAircraftID(userID, id)
	if err != nil {
		return err
	}

	if numberOfFlights > 0 {
		return fmt.Errorf("%w: aircraft has assigned flights", dto.ErrConflict)
	}

	err = a.aircraftRepository.DeleteByUserIDAndID(userID, id)
	if err != nil {
		return err
	}

	return nil
}

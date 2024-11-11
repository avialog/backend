package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -source=logbook.go -destination=logbook_mock.go -package service
type LogbookService interface {
	InsertLogbookEntry(userID string, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	DeleteLogbookEntry(userID string, flightID uint) error
	UpdateLogbookEntry(userID string, flightID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	GetLogbookEntries(userID string, start, end time.Time) ([]dto.LogbookResponse, error)
}

type logbookService struct {
	flightRepository    repository.FlightRepository
	landingRepository   repository.LandingRepository
	passengerRepository repository.PassengerRepository
	aircraftRepository  repository.AircraftRepository
	validator           *validator.Validate
	config              config.Config
}

func newLogbookService(flightRepository repository.FlightRepository, landingRepository repository.LandingRepository,
	passengerRepository repository.PassengerRepository, aircraftRepository repository.AircraftRepository,
	config config.Config, validator *validator.Validate) LogbookService {
	return &logbookService{flightRepository, landingRepository,
		passengerRepository, aircraftRepository, validator, config}
}

func (l *logbookService) InsertLogbookEntry(userID string, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	var logbookResponse dto.LogbookResponse
	landingEntries := make([]dto.LandingEntry, 0)
	passengerEntries := make([]dto.PassengerEntry, 0)

	if _, err := l.aircraftRepository.GetByUserIDAndID(userID, logbookRequest.AircraftID); err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrBadRequest, "aircraft does not belong to user")
		}
		return dto.LogbookResponse{}, err
	}

	tx := l.flightRepository.Begin()

	flight := model.Flight{
		UserID:              userID,
		AircraftID:          logbookRequest.AircraftID,
		TakeoffTime:         logbookRequest.TakeoffTime,
		TakeoffAirportCode:  logbookRequest.TakeoffAirportCode,
		LandingTime:         logbookRequest.LandingTime,
		LandingAirportCode:  logbookRequest.LandingAirportCode,
		Style:               logbookRequest.Style,
		Remarks:             logbookRequest.Remarks,
		PersonalRemarks:     logbookRequest.PersonalRemarks,
		TotalBlockTime:      logbookRequest.TotalBlockTime,
		PilotInCommandTime:  logbookRequest.PilotInCommandTime,
		SecondInCommandTime: logbookRequest.SecondInCommandTime,
		DualReceivedTime:    logbookRequest.DualReceivedTime,
		DualGivenTime:       logbookRequest.DualGivenTime,
		MultiPilotTime:      logbookRequest.MultiPilotTime,
		NightTime:           logbookRequest.NightTime,
		MyRole:              logbookRequest.MyRole,
		IFRTime:             logbookRequest.IFRTime,
		IFRActualTime:       logbookRequest.IFRActualTime,
		IFRSimulatedTime:    logbookRequest.IFRSimulatedTime,
		CrossCountryTime:    logbookRequest.CrossCountryTime,
		SimulatorTime:       logbookRequest.SimulatorTime,
		SignatureURL:        logbookRequest.SignatureURL,
	}

	err := l.validator.Struct(flight)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}

	insertedFlight, err := l.flightRepository.CreateTx(tx, flight)
	if err != nil {
		tx.Rollback()
		return dto.LogbookResponse{}, err
	}

	for _, passengerEntry := range logbookRequest.Passengers {
		passenger := model.Passenger{
			FlightID:     insertedFlight.ID,
			Role:         passengerEntry.Role,
			FirstName:    passengerEntry.FirstName,
			LastName:     passengerEntry.LastName,
			Company:      passengerEntry.Company,
			Phone:        passengerEntry.Phone,
			EmailAddress: passengerEntry.EmailAddress,
			Note:         passengerEntry.Note,
		}

		err := l.validator.Struct(passenger)
		if err != nil {
			tx.Rollback()
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(err, &invalidValidationError) {
				return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				if len(validationErrors) > 0 {
					return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
				}
			}
		}

		passenger, err = l.passengerRepository.CreateTx(tx, passenger)
		if err != nil {
			tx.Rollback()
			return dto.LogbookResponse{}, err
		}
		passengerEntries = append(passengerEntries, dto.PassengerEntry{
			Role:         passenger.Role,
			FirstName:    passenger.FirstName,
			LastName:     passenger.LastName,
			Company:      passenger.Company,
			Phone:        passenger.Phone,
			EmailAddress: passenger.EmailAddress,
			Note:         passenger.Note,
		})
	}

	for _, landingEntry := range logbookRequest.Landings {
		landing := model.Landing{
			FlightID:     insertedFlight.ID,
			ApproachType: landingEntry.ApproachType,
			Count:        landingEntry.Count,
			NightCount:   landingEntry.NightCount,
			DayCount:     landingEntry.DayCount,
			AirportCode:  landingEntry.AirportCode,
		}

		err := l.validator.Struct(landing)
		if err != nil {
			tx.Rollback()
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(err, &invalidValidationError) {
				return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				if len(validationErrors) > 0 {
					return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
				}
			}
		}

		landing, err = l.landingRepository.CreateTx(tx, landing)
		if err != nil {
			tx.Rollback()
			return dto.LogbookResponse{}, err
		}
		landingEntries = append(landingEntries, dto.LandingEntry{
			ApproachType: landing.ApproachType,
			Count:        landing.Count,
			NightCount:   landing.NightCount,
			DayCount:     landing.DayCount,
			AirportCode:  landing.AirportCode,
		})
	}

	logbookResponse = dto.LogbookResponse{
		AircraftID:          insertedFlight.AircraftID,
		TakeoffTime:         insertedFlight.TakeoffTime,
		TakeoffAirportCode:  insertedFlight.TakeoffAirportCode,
		LandingTime:         insertedFlight.LandingTime,
		LandingAirportCode:  insertedFlight.LandingAirportCode,
		Style:               insertedFlight.Style,
		Remarks:             insertedFlight.Remarks,
		PersonalRemarks:     insertedFlight.PersonalRemarks,
		TotalBlockTime:      insertedFlight.TotalBlockTime,
		MyRole:              insertedFlight.MyRole,
		PilotInCommandTime:  insertedFlight.PilotInCommandTime,
		SecondInCommandTime: insertedFlight.SecondInCommandTime,
		DualReceivedTime:    insertedFlight.DualReceivedTime,
		DualGivenTime:       insertedFlight.DualGivenTime,
		MultiPilotTime:      insertedFlight.MultiPilotTime,
		NightTime:           insertedFlight.NightTime,
		IFRTime:             insertedFlight.IFRTime,
		IFRActualTime:       insertedFlight.IFRActualTime,
		IFRSimulatedTime:    insertedFlight.IFRSimulatedTime,
		CrossCountryTime:    insertedFlight.CrossCountryTime,
		SimulatorTime:       insertedFlight.SimulatorTime,
		SignatureURL:        insertedFlight.SignatureURL,
		Passengers:          passengerEntries,
		Landings:            landingEntries,
	}
	if err := tx.Commit().Error; err != nil {
		return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}

	return logbookResponse, nil
}

func (l *logbookService) DeleteLogbookEntry(userID string, flightID uint) error {
	flight, err := l.flightRepository.GetByID(flightID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return fmt.Errorf("%w: %v", dto.ErrBadRequest, "flight not found")
		}
		return err
	}

	if flight.UserID != userID {
		return fmt.Errorf("%w: %v", dto.ErrBadRequest, "flight does not belong to user")
	}

	tx := l.flightRepository.Begin()

	if err := l.landingRepository.DeleteByFlightIDTx(tx, flightID); err != nil {
		tx.Rollback()
		return err
	}

	if err := l.passengerRepository.DeleteByFlightIDTx(tx, flightID); err != nil {
		tx.Rollback()
		return err
	}

	if err := l.flightRepository.DeleteByIDTx(tx, flightID); err != nil {
		tx.Rollback()
		if errors.Is(err, dto.ErrNotFound) {
			return fmt.Errorf("%w: %v", dto.ErrNotFound, "flight not found")
		}

		return fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}
	return nil
}

func (l *logbookService) GetLogbookEntries(userID string, start, end time.Time) ([]dto.LogbookResponse, error) {
	logbookResponses := make([]dto.LogbookResponse, 0)

	flights, err := l.flightRepository.GetByUserIDAndDate(userID, start, end)
	if err != nil {
		return logbookResponses, err
	}

	for _, flight := range flights {
		landings, err := l.landingRepository.GetByFlightID(flight.ID)
		if err != nil {
			return logbookResponses, err
		}

		landingEntries := make([]dto.LandingEntry, 0)
		for _, landing := range landings {
			landingEntries = append(landingEntries, dto.LandingEntry{
				ApproachType: landing.ApproachType,
				Count:        landing.Count,
				NightCount:   landing.NightCount,
				DayCount:     landing.DayCount,
				AirportCode:  landing.AirportCode,
			})
		}

		passengers, err := l.passengerRepository.GetByFlightID(flight.ID)
		if err != nil {
			return logbookResponses, err
		}

		passengerEntries := make([]dto.PassengerEntry, 0)

		for _, passenger := range passengers {
			passengerEntries = append(passengerEntries, dto.PassengerEntry{
				Role:         passenger.Role,
				FirstName:    passenger.FirstName,
				LastName:     passenger.LastName,
				Company:      passenger.Company,
				Phone:        passenger.Phone,
				EmailAddress: passenger.EmailAddress,
				Note:         passenger.Note,
			})
		}

		logbookResponse := dto.LogbookResponse{
			AircraftID:          flight.AircraftID,
			TakeoffTime:         flight.TakeoffTime,
			TakeoffAirportCode:  flight.TakeoffAirportCode,
			LandingTime:         flight.LandingTime,
			LandingAirportCode:  flight.LandingAirportCode,
			Style:               flight.Style,
			Remarks:             flight.Remarks,
			PersonalRemarks:     flight.PersonalRemarks,
			TotalBlockTime:      flight.TotalBlockTime,
			MyRole:              flight.MyRole,
			PilotInCommandTime:  flight.PilotInCommandTime,
			SecondInCommandTime: flight.SecondInCommandTime,
			DualReceivedTime:    flight.DualReceivedTime,
			DualGivenTime:       flight.DualGivenTime,
			MultiPilotTime:      flight.MultiPilotTime,
			NightTime:           flight.NightTime,
			IFRTime:             flight.IFRTime,
			IFRActualTime:       flight.IFRActualTime,
			IFRSimulatedTime:    flight.IFRSimulatedTime,
			CrossCountryTime:    flight.CrossCountryTime,
			SimulatorTime:       flight.SimulatorTime,
			SignatureURL:        flight.SignatureURL,
			Passengers:          passengerEntries,
			Landings:            landingEntries,
		}

		logbookResponses = append(logbookResponses, logbookResponse)
	}

	return logbookResponses, nil
}

func (l *logbookService) UpdateLogbookEntry(userID string, flightID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	var logbookResponse dto.LogbookResponse
	landingEntries := make([]dto.LandingEntry, 0)
	passengerEntries := make([]dto.PassengerEntry, 0)

	flight, err := l.flightRepository.GetByID(flightID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrNotFound, "flight not found")
		}
		return dto.LogbookResponse{}, err
	}

	if flight.UserID != userID {
		return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrBadRequest, "flight does not belong to user")
	}

	if _, err := l.aircraftRepository.GetByUserIDAndID(userID, logbookRequest.AircraftID); err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrBadRequest, "aircraft does not belong to user")
		}
		return dto.LogbookResponse{}, err
	}

	tx := l.flightRepository.Begin()

	flight.AircraftID = logbookRequest.AircraftID
	flight.TakeoffTime = logbookRequest.TakeoffTime
	flight.TakeoffAirportCode = logbookRequest.TakeoffAirportCode
	flight.LandingTime = logbookRequest.LandingTime
	flight.LandingAirportCode = logbookRequest.LandingAirportCode
	flight.Style = logbookRequest.Style
	flight.Remarks = logbookRequest.Remarks
	flight.PersonalRemarks = logbookRequest.PersonalRemarks
	flight.TotalBlockTime = logbookRequest.TotalBlockTime
	flight.PilotInCommandTime = logbookRequest.PilotInCommandTime
	flight.SecondInCommandTime = logbookRequest.SecondInCommandTime
	flight.DualReceivedTime = logbookRequest.DualReceivedTime
	flight.DualGivenTime = logbookRequest.DualGivenTime
	flight.MultiPilotTime = logbookRequest.MultiPilotTime
	flight.NightTime = logbookRequest.NightTime
	flight.MyRole = logbookRequest.MyRole
	flight.IFRTime = logbookRequest.IFRTime
	flight.IFRActualTime = logbookRequest.IFRActualTime
	flight.IFRSimulatedTime = logbookRequest.IFRSimulatedTime
	flight.CrossCountryTime = logbookRequest.CrossCountryTime
	flight.SimulatorTime = logbookRequest.SimulatorTime
	flight.SignatureURL = logbookRequest.SignatureURL

	err = l.validator.Struct(flight)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			if len(validationErrors) > 0 {
				return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
			}
		}
	}

	if _, err := l.flightRepository.SaveTx(tx, flight); err != nil {
		tx.Rollback()
		return dto.LogbookResponse{}, err
	}

	if err := l.passengerRepository.DeleteByFlightIDTx(tx, flightID); err != nil {
		tx.Rollback()
		return dto.LogbookResponse{}, err
	}

	for _, passengerEntry := range logbookRequest.Passengers {
		passenger := model.Passenger{
			FlightID:     flight.ID,
			Role:         passengerEntry.Role,
			FirstName:    passengerEntry.FirstName,
			LastName:     passengerEntry.LastName,
			Company:      passengerEntry.Company,
			Phone:        passengerEntry.Phone,
			EmailAddress: passengerEntry.EmailAddress,
			Note:         passengerEntry.Note,
		}

		err = l.validator.Struct(passenger)
		if err != nil {
			tx.Rollback()
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(err, &invalidValidationError) {
				return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				if len(validationErrors) > 0 {
					return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
				}
			}
		}

		if _, err := l.passengerRepository.CreateTx(tx, passenger); err != nil {
			tx.Rollback()
			return dto.LogbookResponse{}, err
		}
		passengerEntries = append(passengerEntries, dto.PassengerEntry{
			Role:         passenger.Role,
			FirstName:    passenger.FirstName,
			LastName:     passenger.LastName,
			Company:      passenger.Company,
			Phone:        passenger.Phone,
			EmailAddress: passenger.EmailAddress,
			Note:         passenger.Note,
		})
	}

	if err := l.landingRepository.DeleteByFlightIDTx(tx, flightID); err != nil {
		tx.Rollback()
		return dto.LogbookResponse{}, err
	}

	for _, landingEntry := range logbookRequest.Landings {
		landing := model.Landing{
			FlightID:     flight.ID,
			ApproachType: landingEntry.ApproachType,
			Count:        landingEntry.Count,
			NightCount:   landingEntry.NightCount,
			DayCount:     landingEntry.DayCount,
			AirportCode:  landingEntry.AirportCode,
		}

		err = l.validator.Struct(landing)
		if err != nil {
			tx.Rollback()
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(err, &invalidValidationError) {
				return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				if len(validationErrors) > 0 {
					return dto.LogbookResponse{}, fmt.Errorf("%w: invalid data in field: %v", dto.ErrBadRequest, validationErrors[0].Field())
				}
			}
		}

		if _, err := l.landingRepository.CreateTx(tx, landing); err != nil {
			tx.Rollback()
			return dto.LogbookResponse{}, err
		}
		landingEntries = append(landingEntries, dto.LandingEntry{
			ApproachType: landing.ApproachType,
			Count:        landing.Count,
			NightCount:   landing.NightCount,
			DayCount:     landing.DayCount,
			AirportCode:  landing.AirportCode,
		})
	}

	logbookResponse = dto.LogbookResponse{
		AircraftID:          flight.AircraftID,
		TakeoffTime:         flight.TakeoffTime,
		TakeoffAirportCode:  flight.TakeoffAirportCode,
		LandingTime:         flight.LandingTime,
		LandingAirportCode:  flight.LandingAirportCode,
		Style:               flight.Style,
		Remarks:             flight.Remarks,
		PersonalRemarks:     flight.PersonalRemarks,
		TotalBlockTime:      flight.TotalBlockTime,
		MyRole:              flight.MyRole,
		PilotInCommandTime:  flight.PilotInCommandTime,
		SecondInCommandTime: flight.SecondInCommandTime,
		DualReceivedTime:    flight.DualReceivedTime,
		DualGivenTime:       flight.DualGivenTime,
		MultiPilotTime:      flight.MultiPilotTime,
		NightTime:           flight.NightTime,
		IFRTime:             flight.IFRTime,
		IFRActualTime:       flight.IFRActualTime,
		IFRSimulatedTime:    flight.IFRSimulatedTime,
		CrossCountryTime:    flight.CrossCountryTime,
		SimulatorTime:       flight.SimulatorTime,
		SignatureURL:        flight.SignatureURL,
		Passengers:          passengerEntries,
		Landings:            landingEntries,
	}

	if err := tx.Commit().Error; err != nil {
		return dto.LogbookResponse{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, err)
	}

	return logbookResponse, nil
}

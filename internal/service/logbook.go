package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"time"
)

type LogbookService interface {
	InsertLogbookEntry(userID, aircraftID uint, logbookRequest dto.LogbookRequest) error
	DeleteLogbookEntry(userID, id uint) error
	UpdateLogbookEntry(userID, id uint, logbookRequest dto.LogbookRequest) error
	GetLogbookEntries(userID uint, start, end time.Time) ([]dto.LogbookResponse, error)
}

type logbookService struct {
	flightRepository    repository.FlightRepository
	landingRepository   repository.LandingRepository
	passengerRepository repository.PassengerRepository
	config              dto.Config
}

func newLogbookService(flightRepository repository.FlightRepository, landingRepository repository.LandingRepository,
	passengerRepository repository.PassengerRepository, config dto.Config) LogbookService {
	return &logbookService{flightRepository, landingRepository, passengerRepository, config}
}

func (l logbookService) InsertLogbookEntry(userID, aircraftID uint, logbookRequest dto.LogbookRequest) error {

	tx := l.flightRepository.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	flight := model.Flight{
		UserID:              userID,
		AircraftID:          aircraftID,
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
		IFRTime:             logbookRequest.IFRTime,
		IFRActualTime:       logbookRequest.IFRActualTime,
		IFRSimulatedTime:    logbookRequest.IFRSimulatedTime,
		CrossCountryTime:    logbookRequest.CrossCountryTime,
		SimulatorTime:       logbookRequest.SimulatorTime,
		SignatureURL:        logbookRequest.SignatureURL,
	}

	insertedFlight, err := l.flightRepository.SaveTx(tx, flight)
	if err != nil {
		tx.Rollback()
		return err
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

		if _, err := l.passengerRepository.SaveTx(tx, passenger); err != nil {
			tx.Rollback()
			return err
		}
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

		if _, err := l.landingRepository.SaveTx(tx, landing); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (l logbookService) DeleteLogbookEntry(userID, id uint) error {
	// init transaction
	tx := l.flightRepository.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	// transaction operations

	// check if flight exists
	flight, err := l.flightRepository.GetByIDTx(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// check if flight belongs to user
	if flight.UserID != userID {
		tx.Rollback()
		return errors.New("flight does not belong to user")
	}

	// delete all landings
	if err := l.landingRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}
	// delete all passengers
	if err := l.passengerRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}
	// delete flight
	if err := l.flightRepository.DeleteByIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// czy nie zabiera zbyt dużo pamięci/jakieś usprawnienia, paginacja?
func (l logbookService) GetLogbookEntries(userID uint, start, end time.Time) ([]dto.LogbookResponse, error) {
	var logbookResponses []dto.LogbookResponse
	// get all flights for user
	flights, err := l.flightRepository.GetByUserIDAndDate(userID, start, end)
	if err != nil {
		return nil, err
	}

	for _, flight := range flights {

		// get all landings for flight
		landings, err := l.landingRepository.GetByFlightID(flight.ID)
		if err != nil {
			return nil, err
		}

		// create landing entries for flight
		var landingEntries []dto.LandingEntry

		for _, landing := range landings {
			landingEntries = append(landingEntries, dto.LandingEntry{
				ApproachType: landing.ApproachType,
				Count:        landing.Count,
				NightCount:   landing.NightCount,
				DayCount:     landing.DayCount,
				AirportCode:  landing.AirportCode,
			})
		}

		// get all passengers for flight
		passengers, err := l.passengerRepository.GetByFlightID(flight.ID)
		if err != nil {
			return nil, err
		}

		// create passenger entries for flight
		var passengerEntries []dto.PassengerEntry

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

		// create logbook response
		logbookResponse := dto.LogbookResponse{
			TakeoffTime:         flight.TakeoffTime,
			TakeoffAirportCode:  flight.TakeoffAirportCode,
			LandingTime:         flight.LandingTime,
			LandingAirportCode:  flight.LandingAirportCode,
			Style:               flight.Style,
			Remarks:             flight.Remarks,
			PersonalRemarks:     flight.PersonalRemarks,
			TotalBlockTime:      flight.TotalBlockTime,
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

func (l logbookService) UpdateLogbookEntry(userID, id uint, logbookRequest dto.LogbookRequest) error {
	// start transaction
	tx := l.flightRepository.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	// Check if flight exists and belongs to user
	flight, err := l.flightRepository.GetByIDTx(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if flight.UserID != userID {
		tx.Rollback()
		return errors.New("flight does not belong to user")
	}

	// Aktualizacja danych lotu
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
	flight.IFRTime = logbookRequest.IFRTime
	flight.IFRActualTime = logbookRequest.IFRActualTime
	flight.IFRSimulatedTime = logbookRequest.IFRSimulatedTime
	flight.CrossCountryTime = logbookRequest.CrossCountryTime
	flight.SimulatorTime = logbookRequest.SimulatorTime
	flight.SignatureURL = logbookRequest.SignatureURL

	// update flight
	if _, err := l.flightRepository.UpdateTx(tx, flight); err != nil {
		tx.Rollback()
		return err
	}

	// delete old passengers
	if err := l.passengerRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	// insert new passengers
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

		if _, err := l.passengerRepository.SaveTx(tx, passenger); err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete old landings
	if err := l.landingRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	// insert new landings
	for _, landingEntry := range logbookRequest.Landings {
		landing := model.Landing{
			FlightID:     flight.ID,
			ApproachType: landingEntry.ApproachType,
			Count:        landingEntry.Count,
			NightCount:   landingEntry.NightCount,
			DayCount:     landingEntry.DayCount,
			AirportCode:  landingEntry.AirportCode,
		}

		if _, err := l.landingRepository.SaveTx(tx, landing); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

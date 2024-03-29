package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"time"
)

type LogbookService interface {
	InsertLogbookEntry(userID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	DeleteLogbookEntry(userID, id uint) error
	UpdateLogbookEntry(userID, id uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	GetLogbookEntries(userID uint, start, end time.Time) ([]dto.LogbookResponse, error)
}

type logbookService struct {
	flightRepository    repository.FlightRepository
	landingRepository   repository.LandingRepository
	passengerRepository repository.PassengerRepository
	aircraftRepository  repository.AircraftRepository
	config              dto.Config
}

func newLogbookService(flightRepository repository.FlightRepository, landingRepository repository.LandingRepository,
	passengerRepository repository.PassengerRepository, aircraftRepository repository.AircraftRepository, config dto.Config) LogbookService {
	return &logbookService{flightRepository, landingRepository, passengerRepository, aircraftRepository, config}
}

func (l *logbookService) InsertLogbookEntry(userID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	var logbookResponse dto.LogbookResponse
	landingEntries := make([]dto.LandingEntry, 0)
	passengerEntries := make([]dto.PassengerEntry, 0)

	if _, err := l.aircraftRepository.GetByUserIDAndID(userID, logbookRequest.AircraftID); err != nil {
		return dto.LogbookResponse{}, errors.New("aircraft does not belong to user")
	}

	tx := l.flightRepository.Begin()

	defer func() {
		tx.Commit()
	}()

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

		passenger, err := l.passengerRepository.SaveTx(tx, passenger)
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
		landing, err := l.landingRepository.SaveTx(tx, landing)
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
		TakeoffTime:         insertedFlight.TakeoffTime,
		TakeoffAirportCode:  insertedFlight.TakeoffAirportCode,
		LandingTime:         insertedFlight.LandingTime,
		LandingAirportCode:  insertedFlight.LandingAirportCode,
		Style:               insertedFlight.Style,
		Remarks:             insertedFlight.Remarks,
		PersonalRemarks:     insertedFlight.PersonalRemarks,
		TotalBlockTime:      insertedFlight.TotalBlockTime,
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
	return logbookResponse, nil
}

func (l *logbookService) DeleteLogbookEntry(userID, id uint) error {
	flight, err := l.flightRepository.GetByID(id)
	if err != nil {
		return err
	}

	if flight.UserID != userID {
		return errors.New("flight does not belong to user")
	}

	tx := l.flightRepository.Begin()

	defer func() {
		tx.Commit()
	}()

	if err := l.landingRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := l.passengerRepository.DeleteByFlightIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := l.flightRepository.DeleteByIDTx(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (l *logbookService) GetLogbookEntries(userID uint, start, end time.Time) ([]dto.LogbookResponse, error) {
	logbookResponses := make([]dto.LogbookResponse, 0)

	flights, err := l.flightRepository.GetByUserIDAndDate(userID, start, end)
	if err != nil {
		return nil, err
	}

	for _, flight := range flights {

		landings, err := l.landingRepository.GetByFlightID(flight.ID)
		if err != nil {
			return nil, err
		}

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

		passengers, err := l.passengerRepository.GetByFlightID(flight.ID)
		if err != nil {
			return nil, err
		}

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

func (l *logbookService) UpdateLogbookEntry(userID, id uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	var logbookResponse dto.LogbookResponse
	landingEntries := make([]dto.LandingEntry, 0)
	passengerEntries := make([]dto.PassengerEntry, 0)

	flight, err := l.flightRepository.GetByID(id)
	if err != nil {
		return dto.LogbookResponse{}, err
	}

	if flight.UserID != userID {
		return dto.LogbookResponse{}, errors.New("flight does not belong to user")
	}

	if _, err := l.aircraftRepository.GetByUserIDAndID(userID, logbookRequest.AircraftID); err != nil {
		return dto.LogbookResponse{}, errors.New("aircraft does not belong to user")
	}

	tx := l.flightRepository.Begin()

	defer func() {
		tx.Commit()
	}()

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
	flight.IFRTime = logbookRequest.IFRTime
	flight.IFRActualTime = logbookRequest.IFRActualTime
	flight.IFRSimulatedTime = logbookRequest.IFRSimulatedTime
	flight.CrossCountryTime = logbookRequest.CrossCountryTime
	flight.SimulatorTime = logbookRequest.SimulatorTime
	flight.SignatureURL = logbookRequest.SignatureURL

	if _, err := l.flightRepository.UpdateTx(tx, flight); err != nil {
		tx.Rollback()
		return dto.LogbookResponse{}, err
	}

	if err := l.passengerRepository.DeleteByFlightIDTx(tx, id); err != nil {
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

		if _, err := l.passengerRepository.SaveTx(tx, passenger); err != nil {
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

	if err := l.landingRepository.DeleteByFlightIDTx(tx, id); err != nil {
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

		if _, err := l.landingRepository.SaveTx(tx, landing); err != nil {
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

	return logbookResponse, nil
}

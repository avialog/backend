package service

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/go-playground/validator/v10"

	"github.com/avialog/backend/internal/pdfexport"
)

//go:generate mockgen -source=logbook.go -destination=logbook_mock.go -package service
type LogbookService interface {
	InsertLogbookEntry(userID string, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	DeleteLogbookEntry(userID string, flightID uint) error
	UpdateLogbookEntry(userID string, flightID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error)
	GetLogbookEntries(userID string, start, end time.Time) ([]dto.LogbookResponse, error)
	GeneratePDF(userID string) ([]byte, error)
}

type logbookService struct {
	flightRepository    repository.FlightRepository
	landingRepository   repository.LandingRepository
	passengerRepository repository.PassengerRepository
	aircraftRepository  repository.AircraftRepository
	userRepository      repository.UserRepository
	validator           *validator.Validate
	config              config.Config
}

func newLogbookService(flightRepository repository.FlightRepository, landingRepository repository.LandingRepository,
	passengerRepository repository.PassengerRepository, aircraftRepository repository.AircraftRepository,
	userRepository repository.UserRepository, config config.Config, validator *validator.Validate) LogbookService {
	return &logbookService{flightRepository, landingRepository,
		passengerRepository, aircraftRepository, userRepository, validator, config}
}
func (l *logbookService) GeneratePDF(userID string) ([]byte, error) {
	user, err := l.userRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	// Create export configuration
	exportConfig := pdfexport.ExportPDF{
		LogbookRows:          18,
		Fill:                 3,
		LeftMargin:           10.0,
		TopMargin:            30.0,
		BodyRow:              5.0,
		FooterRow:            6.0,
		IncludeSignature:     true,
		TimeFieldsAutoFormat: 1,
		Columns: pdfexport.ColumnsWidth{
			Col1:  12.2, // Date
			Col2:  8.25, // Departure
			Col3:  8.25,
			Col4:  8.25, // Arrival
			Col5:  8.25,
			Col6:  10.0, // Aircraft
			Col7:  12.9,
			Col8:  11.2,  // SE
			Col9:  11.2,  // ME
			Col10: 11.2,  // MCC
			Col11: 11.2,  // Total time
			Col12: 22.86, // PIC name
			Col13: 8.38,  // Landings
			Col14: 8.38,
			Col15: 11.2, // Night
			Col16: 11.2, // IFR
			Col17: 11.2, // PIC
			Col18: 11.2, // COP
			Col19: 11.2, // Dual
			Col20: 11.2, // Instr
			Col21: 11.2, // FSTD
			Col22: 11.2,
			Col23: 33.8, // Remarks
		},
		Headers: pdfexport.ColumnsHeader{
			Date:      "DATE",
			Departure: "DEPARTURE",
			Arrival:   "ARRIVAL",
			Aircraft:  "AIRCRAFT",
			SPT:       "SINGLE PILOT TIME",
			MCC:       "MULTI PILOT TIME",
			Total:     "TOTAL TIME",
			PICName:   "PIC NAME",
			Landings:  "LANDINGS",
			OCT:       "OPERATIONAL CONDITION TIME",
			PFT:       "PILOT FUNCTION TIME",
			FSTD:      "FSTD SESSION",
			Remarks:   "REMARKS AND ENDORSEMENTS",
			DepPlace:  "Place",
			DepTime:   "Time",
			ArrPlace:  "Place",
			ArrTime:   "Time",
			Model:     "Type",
			Reg:       "Reg",
			SE:        "SE",
			ME:        "ME",
			LandDay:   "Day",
			LandNight: "Night",
			Night:     "Night",
			IFR:       "IFR",
			PIC:       "PIC",
			COP:       "COP",
			Dual:      "DUAL",
			Instr:     "INSTR",
			SimType:   "Type",
			SimTime:   "Time",
		},
	}

	exporter, err := pdfexport.NewPDFExporter(
		"A4",
		*user.FirstName+" "+*user.LastName,
		*user.LicenseNumber, // number of license
		*user.Address,
		"",
		"",
		exportConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF exporter: %w", err)
	}

	err = exporter.ExportA4(l.mapToSingleLogbookEntry(userID), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to export PDF: %w", err)
	}

	return buf.Bytes(), nil
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
			FlightID:            flight.ID,
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

func (l *logbookService) mapToSingleLogbookEntry(userID string) []pdfexport.SingleLogbookEntry {
	userFlights, err := l.flightRepository.GetFlightForLogbook(userID)
	if err != nil {
		return nil
	}

	var logbookEntries []pdfexport.SingleLogbookEntry

	for _, flight := range userFlights {
		if flight.Remarks == nil {
			fmt.Printf("Warning: flight %d has nil Remarks\n", flight.ID)
			continue
		}

		// Get date and times from takeoff/landing
		date, depTime := l.transformDateTime(flight.TakeoffTime)
		_, arrTime := l.transformDateTime(flight.LandingTime)

		// Calculate total landings (day/night)
		dayLandings := 0
		nightLandings := 0
		for _, landing := range flight.Landings {
			if landing.DayCount != nil {
				dayLandings += int(*landing.DayCount)
			}
			if landing.NightCount != nil {
				nightLandings += int(*landing.NightCount)
			}
		}

		// Determine PIC name
		picName := ""
		if flight.MyRole == "PIC" {
			picName = "SELF"
		} else {

			for _, passenger := range flight.Passengers {
				if passenger.Role == "PIC" {
					picName = fmt.Sprintf("%s %s", passenger.FirstName, *passenger.LastName)
					break
				}
			}
		}

		isMultiPilot := false
		qualifiedPassengers := 0

		if flight.MyRole == "PIC" || flight.MyRole == "SIC" ||
			flight.MyRole == "INS" || flight.MyRole == "EXM" ||
			flight.MyRole == "P1S" {
			qualifiedPassengers++
		}

		for _, passenger := range flight.Passengers {
			if passenger.Role == "PIC" || passenger.Role == "SIC" ||
				passenger.Role == "INS" || passenger.Role == "EXM" ||
				passenger.Role == "P1S" {
				qualifiedPassengers++
			}
		}

		isMultiPilot = qualifiedPassengers >= 2

		// Convert all duration fields
		var totalTime string
		if flight.TotalBlockTime != nil {
			totalTime = l.translateDurationToString(flight.TotalBlockTime.Nanoseconds())
		}

		// Determine SE/ME/MCC times based on aircraft type and multi-pilot status
		var seTime, meTime, mccTime string
		if !isMultiPilot {
			if flight.Aircraft.IsSingleEngine == "true" {
				seTime = totalTime
			} else {
				meTime = totalTime
			}
		} else {
			mccTime = totalTime
		}

		// Convert other time fields
		var picTime, sicTime, nightTime, ifrTime,
			dualTime, instrTime, simTime string

		if flight.PilotInCommandTime != nil {
			picTime = l.translateDurationToString(flight.PilotInCommandTime.Nanoseconds())
		}
		if flight.SecondInCommandTime != nil {
			sicTime = l.translateDurationToString(flight.SecondInCommandTime.Nanoseconds())
		}
		if flight.NightTime != nil {
			nightTime = l.translateDurationToString(flight.NightTime.Nanoseconds())
		}
		if flight.IFRTime != nil {
			ifrTime = l.translateDurationToString(flight.IFRTime.Nanoseconds())
		}
		if flight.DualReceivedTime != nil {
			dualTime = l.translateDurationToString(flight.DualReceivedTime.Nanoseconds())
		}
		if flight.DualGivenTime != nil {
			instrTime = l.translateDurationToString(flight.DualGivenTime.Nanoseconds())
		}
		if flight.SimulatorTime != nil {
			simTime = l.translateDurationToString(flight.SimulatorTime.Nanoseconds())
		}

		entry := pdfexport.SingleLogbookEntry{
			Date: date,
			Departure: struct {
				Place string
				Time  string
			}{
				Place: flight.TakeoffAirportCode,
				Time:  depTime,
			},
			Arrival: struct {
				Place string
				Time  string
			}{
				Place: flight.LandingAirportCode,
				Time:  arrTime,
			},
			Aircraft: struct {
				Model string
				Reg   string
			}{
				Model: flight.Aircraft.AircraftModel,
				Reg:   flight.Aircraft.RegistrationNumber,
			},
			Time: struct {
				SE         string
				ME         string
				MCC        string
				Total      string
				Night      string
				IFR        string
				PIC        string
				CoPilot    string
				Dual       string
				Instructor string
			}{
				SE:         seTime,
				ME:         meTime,
				MCC:        mccTime,
				Total:      totalTime,
				Night:      nightTime,
				IFR:        ifrTime,
				PIC:        picTime,
				CoPilot:    sicTime,
				Dual:       dualTime,
				Instructor: instrTime,
			},
			Landings: struct {
				Day   int
				Night int
			}{
				Day:   dayLandings,
				Night: nightLandings,
			},
			SIM: struct {
				Type string
				Time string
			}{
				Type: func() string {
					if flight.FSTDtype != nil {
						return *flight.FSTDtype
					}
					return ""
				}(),
				Time: simTime,
			},
			PIC:     picName,
			Remarks: *flight.Remarks,
		}

		logbookEntries = append(logbookEntries, entry)
	}

	return logbookEntries
}

func (l *logbookService) translateDurationToString(durationNano int64) string {
	if durationNano == 0 {
		return ""
	}

	// Convert directly from seconds to hours and minutes
	hours := durationNano / 3600
	minutes := (durationNano % 3600) / 60

	return fmt.Sprintf("%d:%02d", hours, minutes)
}

func (l *logbookService) transformDateTime(t time.Time) (date string, timeStr string) {
	if t.IsZero() {
		return "", ""
	}

	date = t.Format("02.01.2006") // DD.MM.YYYY
	timeStr = t.Format("15:04")   // HH:MM (24-hour format)

	return date, timeStr
}

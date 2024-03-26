package service

import (
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"time"
)

var _ = Describe("LogbookService", func() {
	var (
		logbookService    LogbookService
		flightRepoCtrl    *gomock.Controller
		flightRepoMock    *repository.MockFlightRepository
		landingRepoCtrl   *gomock.Controller
		landingRepoMock   *repository.MockLandingRepository
		passengerRepoCtrl *gomock.Controller
		passengerRepoMock *repository.MockPassengerRepository
		aircraftRepoCtrl  *gomock.Controller
		aircraftRepoMock  *repository.MockAircraftRepository
		databaseCtrl      *gomock.Controller
		databaseMock      *repository.MockDatabase
		validator         *validator.Validate
		logbookRequest    dto.LogbookRequest
		fixedTime         time.Time
		mockFlight        model.Flight
		//mockPassengerOne         model.Passenger
		mockInsertedPassengerOne model.Passenger

		//mockPassengerTwo         model.Passenger
		mockInsertedPassengerTwo model.Passenger
		mockLandingOne           model.Landing
		mockInsertedLandingOne   model.Landing
		mockLandingTwo           model.Landing
		mockInsertedLandingTwo   model.Landing
		mockInsertedFlight       model.Flight
	)

	BeforeEach(func() {
		fixedTime = time.Date(2024, 3, 25, 23, 12, 51, 0, time.UTC)
		flightRepoCtrl = gomock.NewController(GinkgoT())
		flightRepoMock = repository.NewMockFlightRepository(flightRepoCtrl)
		landingRepoCtrl = gomock.NewController(GinkgoT())
		landingRepoMock = repository.NewMockLandingRepository(landingRepoCtrl)
		passengerRepoCtrl = gomock.NewController(GinkgoT())
		passengerRepoMock = repository.NewMockPassengerRepository(passengerRepoCtrl)
		aircraftRepoCtrl = gomock.NewController(GinkgoT())
		aircraftRepoMock = repository.NewMockAircraftRepository(aircraftRepoCtrl)
		databaseCtrl = gomock.NewController(GinkgoT())
		databaseMock = repository.NewMockDatabase(databaseCtrl)
		validator = utils.GetValidator()
		logbookService = newLogbookService(flightRepoMock, landingRepoMock, passengerRepoMock, aircraftRepoMock, dto.Config{}, validator)
		logbookRequest = dto.LogbookRequest{
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             "Remarks",
			PersonalRemarks:     "Personal Remarks",
			TotalBlockTime:      1 * time.Hour,
			PilotInCommandTime:  2 * time.Hour,
			SecondInCommandTime: 3 * time.Hour,
			DualReceivedTime:    4 * time.Hour,
			DualGivenTime:       5 * time.Hour,
			MultiPilotTime:      6 * time.Hour,
			NightTime:           7 * time.Hour,
			IFRTime:             8 * time.Hour,
			IFRActualTime:       9 * time.Hour,
			IFRSimulatedTime:    10 * time.Hour,
			CrossCountryTime:    11 * time.Hour,
			SimulatorTime:       12 * time.Hour,
			SignatureURL:        "https://signature.com",
			Passengers: []dto.PassengerEntry{
				{
					Role:         model.RolePilotInCommand,
					FirstName:    "John",
					LastName:     "Doe",
					Company:      "Company",
					Phone:        "1234567890",
					EmailAddress: "test@test.com",
					Note:         "Note",
				},
				{
					Role:         model.RoleSecondInCommand,
					FirstName:    "Jane",
					LastName:     "Doe",
					Company:      "Company",
					Phone:        "1234567890",
					EmailAddress: "testing@test.com",
					Note:         "Note",
				},
			},
			Landings: []dto.LandingEntry{
				{
					ApproachType: model.ApproachTypeVisual,
					Count:        1,
					NightCount:   2,
					DayCount:     3,
					AirportCode:  "SFO",
				},
				{
					ApproachType: model.ApproachTypeVisual,
					Count:        4,
					NightCount:   5,
					DayCount:     6,
					AirportCode:  "LAX",
				},
			},
		}
		mockFlight = model.Flight{
			UserID:              uint(2),
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             "Remarks",
			PersonalRemarks:     "Personal Remarks",
			TotalBlockTime:      1 * time.Hour,
			PilotInCommandTime:  2 * time.Hour,
			SecondInCommandTime: 3 * time.Hour,
			DualReceivedTime:    4 * time.Hour,
			DualGivenTime:       5 * time.Hour,
			MultiPilotTime:      6 * time.Hour,
			NightTime:           7 * time.Hour,
			IFRTime:             8 * time.Hour,
			IFRActualTime:       9 * time.Hour,
			IFRSimulatedTime:    10 * time.Hour,
			CrossCountryTime:    11 * time.Hour,
			SimulatorTime:       12 * time.Hour,
			SignatureURL:        "https://signature.com",
		}
		mockInsertedFlight = model.Flight{
			Model:               gorm.Model{ID: uint(3)},
			UserID:              uint(2),
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             "Remarks",
			PersonalRemarks:     "Personal Remarks",
			TotalBlockTime:      1 * time.Hour,
			PilotInCommandTime:  2 * time.Hour,
			SecondInCommandTime: 3 * time.Hour,
			DualReceivedTime:    4 * time.Hour,
			DualGivenTime:       5 * time.Hour,
			MultiPilotTime:      6 * time.Hour,
			NightTime:           7 * time.Hour,
			IFRTime:             8 * time.Hour,
			IFRActualTime:       9 * time.Hour,
			IFRSimulatedTime:    10 * time.Hour,
			CrossCountryTime:    11 * time.Hour,
			SimulatorTime:       12 * time.Hour,
			SignatureURL:        "https://signature.com",
		}
		//mockPassengerOne = model.Passenger{
		//	FlightID:     uint(3),
		//	Role:         model.RolePilotInCommand,
		//	FirstName:    "John",
		//	LastName:     "Doe",
		//	Company:      "Company",
		//	Phone:        "1234567890",
		//	EmailAddress: "test@test.com",
		//	Note:         "Note",
		//}
		mockInsertedPassengerOne = model.Passenger{
			Model:        gorm.Model{ID: uint(1)},
			FlightID:     uint(3),
			Role:         model.RolePilotInCommand,
			FirstName:    "John",
			LastName:     "Doe",
			Company:      "Company",
			Phone:        "1234567890",
			EmailAddress: "test@test.com",
			Note:         "Note",
		}
		//mockPassengerTwo = model.Passenger{
		//	Model:        gorm.Model{ID: uint(2)},
		//	FlightID:     uint(3),
		//	Role:         model.RoleSecondInCommand,
		//	FirstName:    "Jane",
		//	LastName:     "Doe",
		//	Company:      "Company",
		//	Phone:        "1234567890",
		//	EmailAddress: "testing@test.com",
		//	Note:         "Note",
		//}
		mockInsertedPassengerTwo = model.Passenger{
			Model:        gorm.Model{ID: uint(2)},
			FlightID:     uint(3),
			Role:         model.RoleSecondInCommand,
			FirstName:    "Jane",
			LastName:     "Doe",
			Company:      "Company",
			Phone:        "1234567890",
			EmailAddress: "testing@test.com",
			Note:         "Note",
		}
		mockLandingOne = model.Landing{
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        1,
			NightCount:   2,
			DayCount:     3,
			AirportCode:  "SFO",
		}
		mockInsertedLandingOne = model.Landing{
			Model:        gorm.Model{ID: uint(1)},
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        1,
			NightCount:   2,
			DayCount:     3,
			AirportCode:  "SFO",
		}
		mockLandingTwo = model.Landing{
			FlightID:     3,
			ApproachType: model.ApproachTypeVisual,
			Count:        4,
			NightCount:   5,
			DayCount:     6,
			AirportCode:  "LAX",
		}
		mockInsertedLandingTwo = model.Landing{
			Model:        gorm.Model{ID: uint(2)},
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        4,
			NightCount:   5,
			DayCount:     6,
			AirportCode:  "LAX",
		}

	})

	AfterEach(func() {
		flightRepoCtrl.Finish()
		landingRepoCtrl.Finish()
		passengerRepoCtrl.Finish()
		aircraftRepoCtrl.Finish()
		databaseCtrl.Finish()
	})

	Describe("InsertLogbookEntry", func() {
		Context("When the logbook entry is valid", func() {
			It("Should return the logbook entry and no error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID(uint(2), uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, gomock.Any()).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, gomock.Any()).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(mockInsertedLandingOne, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingTwo).Return(mockInsertedLandingTwo, nil)

				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Commit()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry(uint(2), logbookRequest)

				// then

				Expect(err).To(BeNil())
				Expect(logbookResponse).ToNot(BeNil())
				Expect(logbookResponse.TakeoffTime).To(Equal(logbookRequest.TakeoffTime))
				Expect(logbookResponse.TakeoffAirportCode).To(Equal(logbookRequest.TakeoffAirportCode))
				Expect(logbookResponse.LandingTime).To(Equal(logbookRequest.LandingTime))
				Expect(logbookResponse.LandingAirportCode).To(Equal(logbookRequest.LandingAirportCode))
				Expect(logbookResponse.Style).To(Equal(logbookRequest.Style))
				Expect(logbookResponse.Remarks).To(Equal(logbookRequest.Remarks))
				Expect(logbookResponse.PersonalRemarks).To(Equal(logbookRequest.PersonalRemarks))
				Expect(logbookResponse.TotalBlockTime).To(Equal(logbookRequest.TotalBlockTime))
				Expect(logbookResponse.PilotInCommandTime).To(Equal(logbookRequest.PilotInCommandTime))
				Expect(logbookResponse.SecondInCommandTime).To(Equal(logbookRequest.SecondInCommandTime))
				Expect(logbookResponse.DualReceivedTime).To(Equal(logbookRequest.DualReceivedTime))
				Expect(logbookResponse.DualGivenTime).To(Equal(logbookRequest.DualGivenTime))
				Expect(logbookResponse.MultiPilotTime).To(Equal(logbookRequest.MultiPilotTime))
				Expect(logbookResponse.NightTime).To(Equal(logbookRequest.NightTime))
				Expect(logbookResponse.IFRTime).To(Equal(logbookRequest.IFRTime))
				Expect(logbookResponse.IFRActualTime).To(Equal(logbookRequest.IFRActualTime))
				Expect(logbookResponse.IFRSimulatedTime).To(Equal(logbookRequest.IFRSimulatedTime))
				Expect(logbookResponse.CrossCountryTime).To(Equal(logbookRequest.CrossCountryTime))
				Expect(logbookResponse.SimulatorTime).To(Equal(logbookRequest.SimulatorTime))
				Expect(logbookResponse.SignatureURL).To(Equal(logbookRequest.SignatureURL))
				//Expect(logbookResponse.Passengers).To(HaveLen(2))
				//Expect(logbookResponse.Passengers).To(Equal(logbookRequest.Passengers))
				//Expect(logbookResponse.Landings).To(HaveLen(2))
				//Expect(logbookResponse.Landings).To(Equal(logbookRequest.Landings))

			})
		})
	})

})

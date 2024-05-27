package service

import (
	"errors"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/infrastructure"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/util"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"time"
)

var _ = Describe("LogbookService", func() {
	var (
		logbookService           LogbookService
		flightRepoCtrl           *gomock.Controller
		flightRepoMock           *repository.MockFlightRepository
		landingRepoCtrl          *gomock.Controller
		landingRepoMock          *repository.MockLandingRepository
		passengerRepoCtrl        *gomock.Controller
		passengerRepoMock        *repository.MockPassengerRepository
		aircraftRepoCtrl         *gomock.Controller
		aircraftRepoMock         *repository.MockAircraftRepository
		databaseCtrl             *gomock.Controller
		databaseMock             *infrastructure.MockDatabase
		validator                *validator.Validate
		logbookRequest           dto.LogbookRequest
		fixedTime                time.Time
		mockFlight               model.Flight
		mockPassengerOne         model.Passenger
		mockInsertedPassengerOne model.Passenger
		mockPassengerTwo         model.Passenger
		mockInsertedPassengerTwo model.Passenger
		mockLandingOne           model.Landing
		mockInsertedLandingOne   model.Landing
		mockLandingTwo           model.Landing
		mockInsertedLandingTwo   model.Landing
		mockInsertedFlight       model.Flight
		mockFlightBeforeUpdate   model.Flight
		startDate                time.Time
		endDate                  time.Time
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
		databaseMock = infrastructure.NewMockDatabase(databaseCtrl)
		validator = util.GetValidator()
		logbookService = newLogbookService(flightRepoMock, landingRepoMock, passengerRepoMock, aircraftRepoMock, config.Config{}, validator)
		logbookRequest = dto.LogbookRequest{
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             util.String("Remarks"),
			PersonalRemarks:     util.String("Personal Remarks"),
			TotalBlockTime:      util.Duration(1 * time.Hour),
			PilotInCommandTime:  util.Duration(2 * time.Hour),
			SecondInCommandTime: util.Duration(3 * time.Hour),
			DualReceivedTime:    util.Duration(4 * time.Hour),
			DualGivenTime:       util.Duration(5 * time.Hour),
			MultiPilotTime:      util.Duration(6 * time.Hour),
			NightTime:           util.Duration(7 * time.Hour),
			IFRTime:             util.Duration(8 * time.Hour),
			IFRActualTime:       util.Duration(9 * time.Hour),
			IFRSimulatedTime:    util.Duration(10 * time.Hour),
			CrossCountryTime:    util.Duration(11 * time.Hour),
			SimulatorTime:       util.Duration(12 * time.Hour),
			SignatureURL:        util.String("https://signature.com"),
			Passengers: []dto.PassengerEntry{
				{
					Role:         model.RolePilotInCommand,
					FirstName:    "John",
					LastName:     util.String("Doe"),
					Company:      util.String("Company"),
					Phone:        util.String("1234567890"),
					EmailAddress: util.String("test@test.com"),
					Note:         util.String("Note"),
				},
				{
					Role:         model.RoleSecondInCommand,
					FirstName:    "Jane",
					LastName:     util.String("Doe"),
					Company:      util.String("Company"),
					Phone:        util.String("1234567890"),
					EmailAddress: util.String("testing@test.com"),
					Note:         util.String("Note"),
				},
			},
			Landings: []dto.LandingEntry{
				{
					ApproachType: model.ApproachTypeVisual,
					Count:        util.Uint(1),
					NightCount:   util.Uint(2),
					DayCount:     util.Uint(3),
					AirportCode:  util.String("SFO"),
				},
				{
					ApproachType: model.ApproachTypeVisual,
					Count:        util.Uint(4),
					NightCount:   util.Uint(5),
					DayCount:     util.Uint(6),
					AirportCode:  util.String("LAX"),
				},
			},
		}
		mockFlight = model.Flight{
			UserID:              "2",
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             util.String("Remarks"),
			PersonalRemarks:     util.String("Personal Remarks"),
			TotalBlockTime:      util.Duration(1 * time.Hour),
			PilotInCommandTime:  util.Duration(2 * time.Hour),
			SecondInCommandTime: util.Duration(3 * time.Hour),
			DualReceivedTime:    util.Duration(4 * time.Hour),
			DualGivenTime:       util.Duration(5 * time.Hour),
			MultiPilotTime:      util.Duration(6 * time.Hour),
			NightTime:           util.Duration(7 * time.Hour),
			IFRTime:             util.Duration(8 * time.Hour),
			IFRActualTime:       util.Duration(9 * time.Hour),
			IFRSimulatedTime:    util.Duration(10 * time.Hour),
			CrossCountryTime:    util.Duration(11 * time.Hour),
			SimulatorTime:       util.Duration(12 * time.Hour),
			SignatureURL:        util.String("https://signature.com"),
		}
		mockInsertedFlight = model.Flight{
			Model:               gorm.Model{ID: uint(3)},
			UserID:              "2",
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             util.String("Remarks"),
			PersonalRemarks:     util.String("Personal Remarks"),
			TotalBlockTime:      util.Duration(1 * time.Hour),
			PilotInCommandTime:  util.Duration(2 * time.Hour),
			SecondInCommandTime: util.Duration(3 * time.Hour),
			DualReceivedTime:    util.Duration(4 * time.Hour),
			DualGivenTime:       util.Duration(5 * time.Hour),
			MultiPilotTime:      util.Duration(6 * time.Hour),
			NightTime:           util.Duration(7 * time.Hour),
			IFRTime:             util.Duration(8 * time.Hour),
			IFRActualTime:       util.Duration(9 * time.Hour),
			IFRSimulatedTime:    util.Duration(10 * time.Hour),
			CrossCountryTime:    util.Duration(11 * time.Hour),
			SimulatorTime:       util.Duration(12 * time.Hour),
			SignatureURL:        util.String("https://signature.com"),
		}
		mockPassengerOne = model.Passenger{
			FlightID:     uint(3),
			Role:         model.RolePilotInCommand,
			FirstName:    "John",
			LastName:     util.String("Doe"),
			Company:      util.String("Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("test@test.com"),
			Note:         util.String("Note"),
		}
		mockInsertedPassengerOne = model.Passenger{
			Model:        gorm.Model{ID: uint(1)},
			FlightID:     uint(3),
			Role:         model.RolePilotInCommand,
			FirstName:    "John",
			LastName:     util.String("Doe"),
			Company:      util.String("Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("test@test.com"),
			Note:         util.String("Note"),
		}
		mockPassengerTwo = model.Passenger{
			FlightID:     uint(3),
			Role:         model.RoleSecondInCommand,
			FirstName:    "Jane",
			LastName:     util.String("Doe"),
			Company:      util.String("Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("testing@test.com"),
			Note:         util.String("Note"),
		}
		mockInsertedPassengerTwo = model.Passenger{
			Model:        gorm.Model{ID: uint(2)},
			FlightID:     uint(3),
			Role:         model.RoleSecondInCommand,
			FirstName:    "Jane",
			LastName:     util.String("Doe"),
			Company:      util.String("Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("testing@test.com"),
			Note:         util.String("Note"),
		}
		mockLandingOne = model.Landing{
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        util.Uint(1),
			NightCount:   util.Uint(2),
			DayCount:     util.Uint(3),
			AirportCode:  util.String("SFO"),
		}
		mockInsertedLandingOne = model.Landing{
			Model:        gorm.Model{ID: uint(1)},
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        util.Uint(1),
			NightCount:   util.Uint(2),
			DayCount:     util.Uint(3),
			AirportCode:  util.String("SFO"),
		}
		mockLandingTwo = model.Landing{
			FlightID:     3,
			ApproachType: model.ApproachTypeVisual,
			Count:        util.Uint(4),
			NightCount:   util.Uint(5),
			DayCount:     util.Uint(6),
			AirportCode:  util.String("LAX"),
		}
		mockInsertedLandingTwo = model.Landing{
			Model:        gorm.Model{ID: uint(2)},
			FlightID:     uint(3),
			ApproachType: model.ApproachTypeVisual,
			Count:        util.Uint(4),
			NightCount:   util.Uint(5),
			DayCount:     util.Uint(6),
			AirportCode:  util.String("LAX"),
		}
		mockFlightBeforeUpdate = model.Flight{
			Model:               gorm.Model{ID: uint(3)},
			UserID:              "2",
			AircraftID:          uint(1),
			TakeoffTime:         fixedTime,
			TakeoffAirportCode:  "SFO",
			LandingTime:         fixedTime,
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             util.String("MRemarks"),
			PersonalRemarks:     util.String("Personal Remarks"),
			TotalBlockTime:      util.Duration(6 * time.Hour),
			PilotInCommandTime:  util.Duration(2 * time.Hour),
			SecondInCommandTime: util.Duration(155 * time.Hour),
			DualReceivedTime:    util.Duration(3 * time.Hour),
			DualGivenTime:       util.Duration(3 * time.Hour),
			MultiPilotTime:      util.Duration(2 * time.Hour),
			NightTime:           util.Duration(5 * time.Hour),
			IFRTime:             util.Duration(7 * time.Hour),
			IFRActualTime:       util.Duration(1 * time.Hour),
			IFRSimulatedTime:    util.Duration(10 * time.Hour),
			CrossCountryTime:    util.Duration(11 * time.Hour),
			SimulatorTime:       util.Duration(6 * time.Hour),
			SignatureURL:        util.String("https://signature.com"),
		}
		startDate = time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(2022, time.March, 28, 0, 0, 0, 0, time.UTC)
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
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(mockInsertedLandingOne, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingTwo).Return(mockInsertedLandingTwo, nil)

				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: nil})

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).To(BeNil())
				Expect(logbookResponse.AircraftID).To(Equal(uint(1)))
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
				Expect(logbookResponse.Passengers).To(HaveLen(2))
				Expect(logbookResponse.Passengers).To(Equal(logbookRequest.Passengers))
				Expect(logbookResponse.Landings).To(HaveLen(2))
				Expect(logbookResponse.Landings).To(Equal(logbookRequest.Landings))

			})
		})
		Context("when commit fails", func() {
			It("Should return an error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(mockInsertedLandingOne, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingTwo).Return(mockInsertedLandingTwo, nil)

				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: errors.New("failed to commit")})

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("internal failure: failed to commit"))
			})
		})
		// it also covers the case when the user does not provide the aircraft id in the request
		Context("When the aircraft does not exist or the user does not own the aircraft", func() {
			It("Should return an error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, errors.New("aircraft not found"))

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("aircraft not found"))
			})
		})
		Context("When flight to insert missing takeoff time", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.TakeoffTime = time.Time{}
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: TakeoffTime"))
			})
		})
		Context("When flight to insert missing takeoff airport code", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.TakeoffAirportCode = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: TakeoffAirportCode"))
			})
		})
		Context("When flight to insert missing landing time", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.LandingTime = time.Time{}
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: LandingTime"))
			})
		})
		Context("When flight to insert missing landing airport code", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.LandingAirportCode = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: LandingAirportCode"))
			})
		})
		Context("When flight to insert missing style", func() {
			It("Should return an error and ", func() {
				// given
				logbookRequest.Style = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Style"))

			})
		})
		Context("when flight to insert have invalid style", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.Style = "invalidStyle"
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Style"))
			})
		})
		Context("when creating flight failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(model.Flight{}, errors.New("failed to create flight"))
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to create flight"))
			})
		})

		Context("when passenger to insert missing role", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].Role = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Role"))
			})
		})
		Context("when passenger to insert have invalid role", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].Role = "invalidRole"
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Role"))
			})
		})
		Context("when passenger to insert missing first name", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].FirstName = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: FirstName"))
			})
		})
		Context("when creating passenger failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(model.Passenger{}, errors.New("failed to create passenger"))
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to create passenger"))
			})
		})
		Context("when landing to insert missing approach type", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Landings[0].ApproachType = ""
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: ApproachType"))
			})
		})
		Context("when landing to insert have invalid approach type", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Landings[0].ApproachType = "invalidApproachType"
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)
				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: ApproachType"))
			})
		})
		Context("when creating landing failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().CreateTx(databaseMock, mockFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(model.Landing{}, errors.New("failed to create landing"))
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.InsertLogbookEntry("2", logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to create landing"))
			})
		})
	})

	Describe("DeleteLogbookEntry", func() {
		Context("when deleting goes well", func() {
			It("Should return no error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				flightRepoMock.EXPECT().DeleteByIDTx(databaseMock, uint(1)).Return(nil)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: nil})

				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).To(BeNil())

			})
		})
		Context("when commit fails", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				flightRepoMock.EXPECT().DeleteByIDTx(databaseMock, uint(1)).Return(nil)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: errors.New("failed to commit")})

				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal failure: failed to commit"))
			})
		})
		Context("when fail to fetch flight", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(model.Flight{}, errors.New("failed to fetch flight"))

				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to fetch flight"))
			})
		})
		Context("when user does not own the flight", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(model.Flight{UserID: "3"}, nil)
				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("bad request: flight does not belong to user"))
			})
		})
		Context("when fail to delete landings", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(errors.New("failed to delete landings"))
				databaseMock.EXPECT().Rollback()
				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to delete landings"))
			})
		})
		Context("when fail to delete passengers", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(errors.New("failed to delete passengers"))
				databaseMock.EXPECT().Rollback()
				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to delete passengers"))
			})
		})
		Context("when fail to delete flight", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(1)).Return(nil)
				flightRepoMock.EXPECT().DeleteByIDTx(databaseMock, uint(1)).Return(errors.New("failed to delete flight"))
				databaseMock.EXPECT().Rollback()

				// when
				err := logbookService.DeleteLogbookEntry("2", uint(1))

				// then
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("internal failure: failed to delete flight"))
			})
		})
	})

	Describe("UpdateLogbookEntry", func() {
		Context("when updating goes well", func() {
			It("Should return no error and updated model response", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: nil})
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(mockInsertedLandingOne, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingTwo).Return(mockInsertedLandingTwo, nil)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)
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
				Expect(logbookResponse.Passengers).To(HaveLen(2))
				Expect(logbookResponse.Passengers).To(Equal(logbookRequest.Passengers))
				Expect(logbookResponse.Landings).To(HaveLen(2))
				Expect(logbookResponse.Landings).To(Equal(logbookRequest.Landings))
			})
		})
		Context("when commit fails", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Commit().Return(&gorm.DB{Error: nil})
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(mockInsertedLandingOne, nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingTwo).Return(mockInsertedLandingTwo, nil)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).To(BeNil())
				Expect(logbookResponse).ToNot(BeNil())

			})
		})
		// this test also covers the case when user does not provide aircraft id in the request
		Context("when getting flight failed", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(model.Flight{}, errors.New("failed to fetch flight"))

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to fetch flight"))
			})
		})
		Context("when user does not own the flight", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(model.Flight{UserID: "3"}, nil)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: flight does not belong to user"))
			})
		})
		Context("when aircraft does not exist or user does not own the aircraft", func() {
			It("Should return an error", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockInsertedFlight, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, errors.New("aircraft does not belong to user"))

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("aircraft does not belong to user"))
			})
		})
		Context("when flight to update missing TakeoffTime", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.TakeoffTime = time.Time{}
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: TakeoffTime"))
			})
		})
		Context("when flight to update missing TakeoffAirportCode", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.TakeoffTime = time.Time{}
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: TakeoffTime"))
			})
		})
		Context("when flight to update missing LandingTime", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.LandingTime = time.Time{}
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: LandingTime"))
			})
		})
		Context("when flight to update missing LandingAirportCode", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.LandingAirportCode = ""
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: LandingAirportCode"))

			})
		})
		Context("when flight to update missing Style", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.Style = ""
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Style"))
			})
		})
		Context("when flight to update invalid Style", func() {
			It("Should return an error", func() {
				// given
				logbookRequest.Style = "invalidStyle"
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Style"))

			})
		})
		Context("when updating flight failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(1)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(model.Flight{}, errors.New("failed to update flight"))
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(1), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to update flight"))
			})
		})
		Context("when deleting passengers failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(errors.New("failed to delete passengers"))
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to delete passengers"))

			})
		})
		Context("when passenger to insert missing role", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].Role = ""
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Role"))

			})
		})
		Context("when passenger to insert have invalid role", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].Role = "invalidRole"
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: Role"))
			})
		})
		Context("when passenger to insert missing first name", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Passengers[0].FirstName = ""
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: FirstName"))
			})
		})
		Context("when creating passenger failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(model.Passenger{}, errors.New("failed to create passenger"))

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to create passenger"))
			})
		})
		Context("when deleting landings failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(errors.New("failed to delete landings"))

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to delete landings"))
			})
		})
		Context("when landing to insert missing approach type", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Landings[0].ApproachType = ""
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: ApproachType"))
			})
		})
		Context("when landing to insert have invalid approach type", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				logbookRequest.Landings[0].ApproachType = "invalidApproachType"
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("bad request: invalid data in field: ApproachType"))
			})
		})
		Context("when creating landing failed", func() {
			It("Should return an error and rollback transaction", func() {
				// given
				flightRepoMock.EXPECT().GetByID(uint(3)).Return(mockFlightBeforeUpdate, nil)
				aircraftRepoMock.EXPECT().GetByUserIDAndID("2", uint(1)).Return(model.Aircraft{}, nil)
				flightRepoMock.EXPECT().SaveTx(databaseMock, mockInsertedFlight).Return(mockInsertedFlight, nil)
				passengerRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				flightRepoMock.EXPECT().Begin().Return(databaseMock)
				databaseMock.EXPECT().Rollback()
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerOne).Return(mockInsertedPassengerOne, nil)
				passengerRepoMock.EXPECT().CreateTx(databaseMock, mockPassengerTwo).Return(mockInsertedPassengerTwo, nil)
				landingRepoMock.EXPECT().DeleteByFlightIDTx(databaseMock, uint(3)).Return(nil)
				landingRepoMock.EXPECT().CreateTx(databaseMock, mockLandingOne).Return(model.Landing{}, errors.New("failed to create landing"))

				// when
				logbookResponse, err := logbookService.UpdateLogbookEntry("2", uint(3), logbookRequest)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponse).To(Equal(dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to create landing"))
			})
		})
	})

	Describe("GetLogbookEntries", func() {
		Context("when flights found", func() {
			It("Should return logbook response array with records", func() {
				// given
				mockFlightArr := []model.Flight{
					{
						Model:               gorm.Model{ID: 1},
						UserID:              "1",
						AircraftID:          1,
						TakeoffTime:         time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 25, 1, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(1),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(1),
						DualReceivedTime:    util.Duration(1),
						DualGivenTime:       util.Duration(1),
						MultiPilotTime:      util.Duration(1),
						NightTime:           util.Duration(1),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(1),
						IFRSimulatedTime:    util.Duration(1),
						CrossCountryTime:    util.Duration(1),
						SimulatorTime:       util.Duration(1),
						SignatureURL:        util.String("DUAL2"),
					},
					{
						Model:               gorm.Model{ID: 2},
						UserID:              "1",
						AircraftID:          3,
						TakeoffTime:         time.Date(2022, time.March, 26, 5, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 26, 6, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(2),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(3),
						DualReceivedTime:    util.Duration(3),
						DualGivenTime:       util.Duration(2),
						MultiPilotTime:      util.Duration(5),
						NightTime:           util.Duration(15),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(2),
						IFRSimulatedTime:    util.Duration(21),
						CrossCountryTime:    util.Duration(2),
						SimulatorTime:       util.Duration(2),
						SignatureURL:        util.String("DUAL1"),
					},
				}
				mockPassengerArr := []model.Passenger{
					{
						Model:     gorm.Model{ID: 1},
						FlightID:  1,
						FirstName: "F1",
						LastName:  util.String("L1"),
						Role:      model.RoleSecondInCommand,
					},
				}
				mockPassengerArr2 := []model.Passenger{
					{
						Model:     gorm.Model{ID: 2},
						FlightID:  2,
						FirstName: "F2",
						LastName:  util.String("L2"),
						Role:      model.RoleSecondInCommand,
					},
					{
						Model:     gorm.Model{ID: 3},
						FlightID:  2,
						FirstName: "F3",
						LastName:  util.String("L3"),
						Role:      model.RoleSecondInCommand,
					},
				}
				mockLandingArr := []model.Landing{
					{
						Model:        gorm.Model{ID: 1},
						FlightID:     1,
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(1),
						NightCount:   util.Uint(1),
						DayCount:     util.Uint(4),
						AirportCode:  util.String("L1"),
					},
					{
						Model:        gorm.Model{ID: 2},
						FlightID:     1,
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(1),
						NightCount:   util.Uint(1),
						DayCount:     util.Uint(6),
						AirportCode:  util.String("L1"),
					},
				}
				mockLandingArr2 := make([]model.Landing, 0)
				expectedLogbookResponse := []dto.LogbookResponse{
					{
						AircraftID:          uint(1),
						TakeoffTime:         time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 25, 1, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(1),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(1),
						DualReceivedTime:    util.Duration(1),
						DualGivenTime:       util.Duration(1),
						MultiPilotTime:      util.Duration(1),
						NightTime:           util.Duration(1),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(1),
						IFRSimulatedTime:    util.Duration(1),
						CrossCountryTime:    util.Duration(1),
						SimulatorTime:       util.Duration(1),
						SignatureURL:        util.String("DUAL2"),
						Passengers: []dto.PassengerEntry{
							{
								FirstName: "F1",
								LastName:  util.String("L1"),
								Role:      model.RoleSecondInCommand,
							},
						},
						Landings: []dto.LandingEntry{
							{
								ApproachType: model.ApproachTypeVisual,
								Count:        util.Uint(1),
								NightCount:   util.Uint(1),
								DayCount:     util.Uint(4),
								AirportCode:  util.String("L1"),
							},
							{
								ApproachType: model.ApproachTypeVisual,
								Count:        util.Uint(1),
								NightCount:   util.Uint(1),
								DayCount:     util.Uint(6),
								AirportCode:  util.String("L1"),
							},
						},
					},
					{
						AircraftID:          uint(3),
						TakeoffTime:         time.Date(2022, time.March, 26, 5, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 26, 6, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(2),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(3),
						DualReceivedTime:    util.Duration(3),
						DualGivenTime:       util.Duration(2),
						MultiPilotTime:      util.Duration(5),
						NightTime:           util.Duration(15),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(2),
						IFRSimulatedTime:    util.Duration(21),
						CrossCountryTime:    util.Duration(2),
						SimulatorTime:       util.Duration(2),
						SignatureURL:        util.String("DUAL1"),
						Passengers: []dto.PassengerEntry{
							{
								FirstName: "F2",
								LastName:  util.String("L2"),
								Role:      model.RoleSecondInCommand,
							},
							{
								FirstName: "F3",
								LastName:  util.String("L3"),
								Role:      model.RoleSecondInCommand,
							},
						},
						Landings: []dto.LandingEntry{},
					},
				}

				flightRepoMock.EXPECT().GetByUserIDAndDate("1", startDate, endDate).Return(mockFlightArr, nil)
				landingRepoMock.EXPECT().GetByFlightID(uint(1)).Return(mockLandingArr, nil)
				landingRepoMock.EXPECT().GetByFlightID(uint(2)).Return(mockLandingArr2, nil)
				passengerRepoMock.EXPECT().GetByFlightID(uint(1)).Return(mockPassengerArr, nil)
				passengerRepoMock.EXPECT().GetByFlightID(uint(2)).Return(mockPassengerArr2, nil)

				// when
				logbookResponses, err := logbookService.GetLogbookEntries("1", startDate, endDate)

				// then
				Expect(err).To(BeNil())
				Expect(logbookResponses).To(Equal(expectedLogbookResponse))
			})
		})
		Context("when flights not found", func() {
			It("Should return empty logbook response array", func() {
				// given

				flightRepoMock.EXPECT().GetByUserIDAndDate("1", startDate, endDate).Return([]model.Flight{}, nil)

				// when
				logbookResponses, err := logbookService.GetLogbookEntries("1", startDate, endDate)

				// then
				Expect(err).To(BeNil())
				Expect(logbookResponses).To(HaveLen(0))
				Expect(logbookResponses).To(Equal([]dto.LogbookResponse{}))
			})
		})
		Context("when getting flights failed", func() {
			It("Should return an error and empty array", func() {
				// given
				flightRepoMock.EXPECT().GetByUserIDAndDate("1", startDate, endDate).Return([]model.Flight{}, errors.New("failed to fetch flights"))

				// when
				logbookResponses, err := logbookService.GetLogbookEntries("1", startDate, endDate)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponses).To(HaveLen(0))
				Expect(logbookResponses).To(Equal([]dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to fetch flights"))
			})
		})
		Context("when getting landings failed", func() {
			It("Should return an error and empty array", func() {
				// given
				flightRepoMock.EXPECT().GetByUserIDAndDate("1", startDate, endDate).Return([]model.Flight{
					{
						Model:               gorm.Model{ID: 1},
						UserID:              "1",
						AircraftID:          1,
						TakeoffTime:         time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 25, 1, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(1),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(1),
						DualReceivedTime:    util.Duration(1),
						DualGivenTime:       util.Duration(1),
						MultiPilotTime:      util.Duration(1),
						NightTime:           util.Duration(1),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(1),
						IFRSimulatedTime:    util.Duration(1),
						CrossCountryTime:    util.Duration(1),
						SimulatorTime:       util.Duration(1),
						SignatureURL:        util.String("DUAL2"),
					},
				}, nil)
				landingRepoMock.EXPECT().GetByFlightID(uint(1)).Return([]model.Landing{}, errors.New("failed to fetch landings"))

				// when
				logbookResponses, err := logbookService.GetLogbookEntries("1", startDate, endDate)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponses).To(HaveLen(0))
				Expect(logbookResponses).To(Equal([]dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to fetch landings"))
			})
		})
		Context("when getting passengers failed", func() {
			It("Should return an error and empty array", func() {
				// given
				flightRepoMock.EXPECT().GetByUserIDAndDate("1", startDate, endDate).Return([]model.Flight{
					{
						Model:               gorm.Model{ID: 1},
						UserID:              "1",
						AircraftID:          1,
						TakeoffTime:         time.Date(2022, time.March, 25, 0, 0, 0, 0, time.UTC),
						LandingTime:         time.Date(2022, time.March, 25, 1, 0, 0, 0, time.UTC),
						TakeoffAirportCode:  "T1",
						LandingAirportCode:  "L1",
						Style:               "S1",
						Remarks:             util.String("R1"),
						PersonalRemarks:     util.String("PR1"),
						TotalBlockTime:      util.Duration(1),
						PilotInCommandTime:  util.Duration(1),
						SecondInCommandTime: util.Duration(1),
						DualReceivedTime:    util.Duration(1),
						DualGivenTime:       util.Duration(1),
						MultiPilotTime:      util.Duration(1),
						NightTime:           util.Duration(1),
						IFRTime:             util.Duration(1),
						IFRActualTime:       util.Duration(1),
						IFRSimulatedTime:    util.Duration(1),
						CrossCountryTime:    util.Duration(1),
						SimulatorTime:       util.Duration(1),
						SignatureURL:        util.String("DUAL2"),
					},
				}, nil)
				landingRepoMock.EXPECT().GetByFlightID(uint(1)).Return([]model.Landing{}, nil)
				passengerRepoMock.EXPECT().GetByFlightID(uint(1)).Return([]model.Passenger{}, errors.New("failed to fetch passengers"))

				// when
				logbookResponses, err := logbookService.GetLogbookEntries("1", startDate, endDate)

				// then
				Expect(err).ToNot(BeNil())
				Expect(logbookResponses).To(HaveLen(0))
				Expect(logbookResponses).To(Equal([]dto.LogbookResponse{}))
				Expect(err.Error()).To(Equal("failed to fetch passengers"))
			})
		})
	})
})

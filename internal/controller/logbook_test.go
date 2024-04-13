package controller

import (
	"bytes"
	"encoding/json"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"time"
)

var _ = Describe("LogbookController", func() {
	var (
		logbookController      LogbookController
		logbookServiceCtrl     *gomock.Controller
		logbookServiceMock     *service.MockLogbookService
		w                      *httptest.ResponseRecorder
		ctx                    *gin.Context
		logbookEntriesMock     []dto.LogbookResponse
		expectedLogbookRequest dto.LogbookRequest
		getLogbookRequest      dto.GetLogbookRequest
	)

	BeforeEach(func() {
		logbookServiceCtrl = gomock.NewController(GinkgoT())
		logbookServiceMock = service.NewMockLogbookService(logbookServiceCtrl)
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
		logbookController = newLogbookController(logbookServiceMock)
		logbookEntriesMock = []dto.LogbookResponse{
			{
				TakeoffTime:         time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				TakeoffAirportCode:  "JFK",
				LandingTime:         time.Date(2024, 3, 1, 1, 0, 0, 0, time.UTC),
				LandingAirportCode:  "LAX",
				Style:               model.StyleY,
				Remarks:             util.String("Remarks"),
				PersonalRemarks:     util.String("PersonalRemarks"),
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
				SignatureURL:        util.String("SignatureURL"),
				Passengers: []dto.PassengerEntry{
					{
						Role:         model.RolePilotInCommand,
						FirstName:    "John",
						LastName:     util.String("Doe"),
						Company:      util.String("Company"),
						Phone:        util.String("Phone"),
						EmailAddress: util.String("Email"),
						Note:         util.String("Note"),
					},
					{
						Role:         model.RoleSecondInCommand,
						FirstName:    "Jane",
						LastName:     util.String("Smiths"),
						Company:      util.String("CDA"),
						Phone:        util.String("Smartphone"),
						EmailAddress: util.String("test@test.com"),
						Note:         util.String("Notes"),
					},
				},
				Landings: []dto.LandingEntry{
					{
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(1),
						NightCount:   util.Uint(2),
						DayCount:     util.Uint(3),
						AirportCode:  util.String("JFK"),
					},
					{
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(4),
						NightCount:   util.Uint(5),
						DayCount:     util.Uint(6),
						AirportCode:  util.String("LAX"),
					},
				},
			},
			{
				TakeoffTime:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				TakeoffAirportCode:  "ABC",
				LandingTime:         time.Date(2024, 2, 1, 1, 0, 0, 0, time.UTC),
				LandingAirportCode:  "XYZ",
				Style:               model.StyleZ2,
				Remarks:             util.String("Remarks2"),
				PersonalRemarks:     util.String("PersonalRemarks2"),
				TotalBlockTime:      util.Duration(12 * time.Hour),
				PilotInCommandTime:  util.Duration(22 * time.Hour),
				SecondInCommandTime: util.Duration(32 * time.Hour),
				DualReceivedTime:    util.Duration(42 * time.Hour),
				DualGivenTime:       util.Duration(52 * time.Hour),
				MultiPilotTime:      util.Duration(62 * time.Hour),
				NightTime:           util.Duration(72 * time.Hour),
				IFRTime:             util.Duration(82 * time.Hour),
				IFRActualTime:       util.Duration(92 * time.Hour),
				IFRSimulatedTime:    util.Duration(11 * time.Hour),
				CrossCountryTime:    util.Duration(12 * time.Hour),
				SimulatorTime:       util.Duration(13 * time.Hour),
				SignatureURL:        util.String("SignatureURL2"),
				Passengers: []dto.PassengerEntry{
					{
						Role:         model.RolePilotInCommand,
						FirstName:    "John2",
						LastName:     util.String("Doe2"),
						Company:      util.String("Company2"),
						Phone:        util.String("Phone2"),
						EmailAddress: util.String("Email2"),
						Note:         util.String("Note2"),
					},
					{
						Role:         model.RoleSecondInCommand,
						FirstName:    "Jane2",
						LastName:     util.String("Smiths2"),
						Company:      util.String("CDA2"),
						Phone:        util.String("Smartphone2"),
						EmailAddress: util.String("test@test.com2"),
						Note:         util.String("Notes2"),
					},
				},
				Landings: []dto.LandingEntry{
					{
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(2),
						NightCount:   util.Uint(3),
						DayCount:     util.Uint(4),
						AirportCode:  util.String("JFK2"),
					},
					{
						ApproachType: model.ApproachTypeVisual,
						Count:        util.Uint(5),
						NightCount:   util.Uint(6),
						DayCount:     util.Uint(7),
						AirportCode:  util.String("LAX2"),
					},
				},
			},
		}
		expectedLogbookRequest = dto.LogbookRequest{
			AircraftID:          1,
			TakeoffTime:         time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			TakeoffAirportCode:  "JFK",
			LandingTime:         time.Date(2024, 3, 1, 1, 0, 0, 0, time.UTC),
			LandingAirportCode:  "LAX",
			Style:               model.StyleY,
			Remarks:             util.String("Remarks"),
			PersonalRemarks:     util.String("PersonalRemarks"),
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
			SignatureURL:        util.String("SignatureURL"),
			Passengers: []dto.PassengerEntry{
				{
					Role:         model.RolePilotInCommand,
					FirstName:    "John",
					LastName:     util.String("Doe"),
					Company:      util.String("Company"),
					Phone:        util.String("Phone"),
					EmailAddress: util.String("Email"),
					Note:         util.String("Note"),
				},
				{
					Role:         model.RoleSecondInCommand,
					FirstName:    "Jane",
					LastName:     util.String("Smiths"),
					Company:      util.String("CDA"),
					Phone:        util.String("Smartphone"),
					EmailAddress: util.String("test@test.com"),
					Note:         util.String("Notes"),
				},
			},
			Landings: []dto.LandingEntry{
				{
					ApproachType: model.ApproachTypeVisual,
					Count:        util.Uint(1),
					NightCount:   util.Uint(2),
					DayCount:     util.Uint(3),
					AirportCode:  util.String("JFK"),
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
		getLogbookRequest = dto.GetLogbookRequest{
			Start: util.Int64(time.Now().AddDate(0, 0, -90).Unix()),
			End:   nil,
		}

	})

	AfterEach(func() {
		logbookServiceCtrl.Finish()
	})

	Describe("GetLogbookEntries", func() {
		Context("When the user doesn't send a request and no error occurs.", func() {
			It("should return 200 and only those entries that match the date", func() {
				// given
				expectedLogbookEntriesJSON, err := json.Marshal(logbookEntriesMock)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = httptest.NewRequest("GET", "/logbook", bytes.NewBuffer([]byte("{}")))
				ctx.Set("userID", "1")
				logbookServiceMock.EXPECT().GetLogbookEntries("1", gomock.Any(), gomock.Any()).DoAndReturn(
					func(userID string, start time.Time, end time.Time) ([]dto.LogbookResponse, error) {
						Expect(start).To(BeTemporally("~", time.Now().AddDate(0, 0, -90), time.Second))
						Expect(end).To(BeTemporally("~", time.Now(), time.Second))
						return logbookEntriesMock, nil
					})

				// when
				logbookController.GetLogbookEntries(ctx)

				// then
				Expect(w.Code).To(Equal(200))
				Expect(w.Body).To(MatchJSON(expectedLogbookEntriesJSON))
			})
		})
		Context("When user request fail to bind", func() {
			It("should return 400 and error message", func() {
				// given
				ctx.Request = httptest.NewRequest("GET", "/logbook", bytes.NewBuffer([]byte("")))
				ctx.Set("userID", "1")
				// when
				logbookController.GetLogbookEntries(ctx)

				// then
				Expect(w.Code).To(Equal(400))
				Expect(w.Body).To(MatchJSON(`{"code": 400, "message":"EOF"}`))
			})
		})
		Context("When the user doesn't provide a start date but provide an end date.", func() {
			It("should return 400 and error message", func() {
				// given
				getLogbookRequestJSON, err := json.Marshal(getLogbookRequest)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = httptest.NewRequest("GET", "/logbook", bytes.NewBuffer(getLogbookRequestJSON))
				ctx.Set("userID", "1")

				// when
				logbookController.GetLogbookEntries(ctx)

				// then
				Expect(w.Code).To(Equal(400))
				Expect(w.Body).To(MatchJSON(`{"code": 400, "message":"both start and end time must be provided or neither"}`))
			})
		})
		//Context("When the user doesn't provide an end date but provide a start date.", func() {
		//
		//})

	})
	Describe("InsertLogbookEntry", func() {
		Context("When the user sends a request and no error occurs.", func() {
			It("should return 200 and the inserted entry", func() {
				// given
				expectedLogbookResponseJSON, err := json.Marshal(logbookEntriesMock[0])
				Expect(err).ToNot(HaveOccurred())
				expectedLogbookInsertRequestJSON, err := json.Marshal(expectedLogbookRequest)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = httptest.NewRequest("POST", "/logbook", bytes.NewBuffer(expectedLogbookInsertRequestJSON))
				ctx.Set("userID", "1")
				logbookServiceMock.EXPECT().InsertLogbookEntry("1", expectedLogbookRequest).Return(logbookEntriesMock[0], nil)

				// when
				logbookController.InsertLogbookEntry(ctx)

				// then
				Expect(w.Code).To(Equal(201))
				Expect(w.Body).To(MatchJSON(expectedLogbookResponseJSON))
			})
		})
	})
	Describe("UpdateLogbookEntry", func() {
		Context("When the user sends a request and no error occurs.", func() {
			It("should return 200 and the updated entry", func() {
				// given
				expectedLogbookResponseJSON, err := json.Marshal(logbookEntriesMock[0])
				Expect(err).ToNot(HaveOccurred())
				expectedLogbookUpdateRequestJSON, err := json.Marshal(expectedLogbookRequest)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = httptest.NewRequest("PUT", "/logbook", bytes.NewBuffer(expectedLogbookUpdateRequestJSON))
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})
				logbookServiceMock.EXPECT().UpdateLogbookEntry("1", uint(1), expectedLogbookRequest).Return(logbookEntriesMock[0], nil)

				// when
				logbookController.UpdateLogbookEntry(ctx)

				// then
				Expect(w.Code).To(Equal(200))
				Expect(w.Body).To(MatchJSON(expectedLogbookResponseJSON))
			})
		})
	})
	Describe("DeleteLogbookEntry", func() {
		Context("When the user sends a request and no error occurs.", func() {
			It("should return 200 and message", func() {
				// given

				ctx.Request = httptest.NewRequest("DELETE", "/logbook", nil)
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})
				logbookServiceMock.EXPECT().DeleteLogbookEntry("1", uint(1)).Return(nil)
				// when
				logbookController.DeleteLogbookEntry(ctx)

				// then
				Expect(w.Code).To(Equal(200))
				Expect(w.Body).To(MatchJSON(`{"message":"Logbook entry deleted successfully"}`))
			})
		})
	})
})

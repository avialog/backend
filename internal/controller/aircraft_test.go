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
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("AircraftController", func() {
	var (
		aircraftController     AircraftController
		aircraftServiceCtrl    *gomock.Controller
		aircraftServiceMock    *service.MockAircraftService
		w                      *httptest.ResponseRecorder
		ctx                    *gin.Context
		aircraftArr            []model.Aircraft
		expectedServerResponse []dto.AircraftResponse
	)

	BeforeEach(func() {
		aircraftServiceCtrl = gomock.NewController(GinkgoT())
		aircraftServiceMock = service.NewMockAircraftService(aircraftServiceCtrl)
		aircraftController = newAircraftController(aircraftServiceMock)
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
		aircraftArr = []model.Aircraft{
			{
				Model:              gorm.Model{ID: 1},
				UserID:             "1",
				RegistrationNumber: "PK-ABC",
				AircraftModel:      "Airbus A320",
				Remarks:            nil,
				ImageURL:           util.String("TestURL"),
			},
			{
				Model:              gorm.Model{ID: 2},
				UserID:             "1",
				RegistrationNumber: "PK-DEF",
				AircraftModel:      "Boeing 737",
				Remarks:            util.String("Test Remarks"),
				ImageURL:           nil,
			},
		}
		expectedServerResponse = []dto.AircraftResponse{
			{
				RegistrationNumber: "PK-ABC",
				AircraftModel:      "Airbus A320",
				Remarks:            nil,
				ImageURL:           util.String("TestURL"),
			},
			{
				RegistrationNumber: "PK-DEF",
				AircraftModel:      "Boeing 737",
				Remarks:            util.String("Test Remarks"),
				ImageURL:           nil,
			},
		}
	})

	AfterEach(func() {
		aircraftServiceCtrl.Finish()
	})

	Describe("GetAircraft", func() {
		Context("When the user has aircraft and everything goes well", func() {
			It("Should return 200 and array of aircraft", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(expectedServerResponse)
				Expect(err).NotTo(HaveOccurred())

				ctx.Set("userID", "1")

				aircraftServiceMock.EXPECT().GetUserAircraft("1").Return(aircraftArr, nil)

				// when
				aircraftController.GetAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))

			})
		})
	})

	Describe("InsertAircraft", func() {
		Context("When everything goes well", func() {
			It("Should return 201 and the aircraft", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())
				expectedRequestJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer(expectedRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				aircraftServiceMock.EXPECT().InsertAircraft("1", gomock.Any()).Return(aircraftArr[0], nil)

				// when
				aircraftController.InsertAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))
			})
		})
	})

	Describe("UpdateAircraft", func() {
		Context("When aircraft exists and everything goes well", func() {
			It("Should return 200 and the updated aircraft", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())
				expectedRequestJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer(expectedRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().UpdateAircraft("1", uint(1), gomock.Any()).Return(aircraftArr[0], nil)

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))
			})
		})
	})

	Describe("DeleteAircraft", func() {
		Context("When aircraft exists and everything goes well", func() {
			It("Should return 200 and message", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(gin.H{"message": "Aircraft deleted successfully"})
				Expect(err).NotTo(HaveOccurred())

				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().DeleteAircraft("1", uint(1)).Return(nil)

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))
			})
		})
	})
})


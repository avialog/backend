package controller

import (
	"bytes"
	"encoding/json"
	"errors"
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
		aircraftRequest        dto.AircraftRequest
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
		aircraftRequest = dto.AircraftRequest{
			RegistrationNumber: "PK-ABC",
			AircraftModel:      "Airbus A320",
			Remarks:            nil,
			ImageURL:           util.String("TestURL"),
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
		Context("When error occurred while getting aircraft", func() {
			It("Should return 500 and error message", func() {
				// given
				ctx.Set("userID", "1")
				aircraftServiceMock.EXPECT().GetUserAircraft("1").Return(nil, errors.New("test"))

				// when
				aircraftController.GetAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"test"}`))
			})
		})
	})

	Describe("InsertAircraft", func() {
		Context("When everything goes well", func() {
			It("Should return 201 and the aircraft", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())
				aircraftRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer(aircraftRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				aircraftServiceMock.EXPECT().InsertAircraft("1", aircraftRequest).Return(aircraftArr[0], nil)

				// when
				aircraftController.InsertAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))
			})
		})
		Context("when binding request failed", func() {
			It("Should return 400 and error message", func() {
				// given
				req, err := http.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer([]byte("")))
				Expect(err).NotTo(HaveOccurred())
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				// when
				aircraftController.InsertAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"EOF"}`))
			})
		})
		Context("When bad request error occurred while inserting aircraft", func() {
			It("Should return 400 and error message", func() {
				// given
				aircraftRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer(aircraftRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				aircraftServiceMock.EXPECT().InsertAircraft("1", aircraftRequest).Return(model.Aircraft{}, dto.ErrBadRequest)

				// when
				aircraftController.InsertAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"bad request"}`))
			})
		})
		Context("When internal error occured while inserting aircraft", func() {
			It("Should return 500 and error message", func() {
				// given
				aircraftRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/aircraft", bytes.NewBuffer(aircraftRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				aircraftServiceMock.EXPECT().InsertAircraft("1", aircraftRequest).Return(model.Aircraft{}, dto.ErrInternalFailure)

				// when
				aircraftController.InsertAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure"}`))
			})
		})
	})

	Describe("UpdateAircraft", func() {
		Context("When aircraft exists and everything goes well", func() {
			It("Should return 200 and the updated aircraft", func() {
				// given
				expectedServerResponseJSON, err := json.Marshal(expectedServerResponse[0])
				Expect(err).NotTo(HaveOccurred())
				expectedRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer(expectedRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().UpdateAircraft("1", uint(1), aircraftRequest).Return(aircraftArr[0], nil)

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(expectedServerResponseJSON))
			})
		})
		Context("when parse id failed", func() {
			It("Should return 400 and error message", func() {
				// given
				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer([]byte("")))
				Expect(err).NotTo(HaveOccurred())
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`))
			})
		})
		Context("When binding request failed", func() {
			It("Should return 400 and error message", func() {
				// given
				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer([]byte("")))
				Expect(err).NotTo(HaveOccurred())
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"EOF"}`))
			})
		})
		Context("When bad request error occurred while updating aircraft", func() {
			It("Should return 400 and error message", func() {
				// given
				expectedRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer(expectedRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().UpdateAircraft("1", uint(1), aircraftRequest).Return(model.Aircraft{}, dto.ErrBadRequest)

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"bad request"}`))
			})
		})
		Context("When internal error occured while updating aircraft", func() {
			It("Should return 500 and error message", func() {
				// given
				expectedRequestJSON, err := json.Marshal(aircraftRequest)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/aircraft/1", bytes.NewBuffer(expectedRequestJSON))
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().UpdateAircraft("1", uint(1), aircraftRequest).Return(model.Aircraft{}, dto.ErrInternalFailure)

				// when
				aircraftController.UpdateAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure"}`))
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
		Context("When parse id failed", func() {
			It("Should return 400 and error message", func() {
				// given
				ctx.Set("userID", "1")

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`))
			})
		})
		Context("when delete aircraft return bad request error", func() {
			It("Should return 400 and error message", func() {
				// given
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().DeleteAircraft("1", uint(1)).Return(dto.ErrBadRequest)

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"bad request"}`))

			})
		})
		Context("when delete aircraft return conflict error", func() {
			It("Should return 409 and error message", func() {
				// given
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().DeleteAircraft("1", uint(1)).Return(dto.ErrConflict)

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusConflict))
				Expect(w.Body.String()).To(Equal(`{"code":409,"message":"conflict"}`))
			})
		})
		Context("when delete aircraft return not found error", func() {
			It("Should return 404 and error message", func() {
				// given
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().DeleteAircraft("1", uint(1)).Return(dto.ErrNotFound)

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusNotFound))
				Expect(w.Body.String()).To(Equal(`{"code":404,"message":"not found"}`))
			})
		})
		Context("When internal error occured while deleting aircraft", func() {
			It("Should return 500 and error message", func() {
				// given
				ctx.Set("userID", "1")
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "1"})

				aircraftServiceMock.EXPECT().DeleteAircraft("1", uint(1)).Return(dto.ErrInternalFailure)

				// when
				aircraftController.DeleteAircraft(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure"}`))
			})
		})
	})
})

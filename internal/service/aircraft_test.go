package service

import (
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("AircraftService", func() {
	var (
		aircraftService  AircraftService
		aircraftRepoCtrl *gomock.Controller
		aircraftRepoMock *repository.MockAircraftRepository
		flightRepoCtrl   *gomock.Controller
		flightRepoMock   *repository.MockFlightRepository
		aircraftRequest  dto.AircraftRequest
		mockAircraft     model.Aircraft
		mockAircraftArr  []model.Aircraft
	)

	BeforeEach(func() {
		aircraftRepoCtrl = gomock.NewController(GinkgoT())
		aircraftRepoMock = repository.NewMockAircraftRepository(aircraftRepoCtrl)
		flightRepoCtrl = gomock.NewController(GinkgoT())
		flightRepoMock = repository.NewMockFlightRepository(flightRepoCtrl)
		aircraftService = newAircraftService(aircraftRepoMock, flightRepoMock, dto.Config{})
		aircraftRequest = dto.AircraftRequest{
			AircraftModel:      "Cessna 172",
			RegistrationNumber: "N12345",
			ImageURL:           "https://example.com/image.jpg",
			Remarks:            "This is a test aircraft",
		}
		mockAircraft = model.Aircraft{
			UserID:             uint(1),
			AircraftModel:      "Cessna 172",
			RegistrationNumber: "N12345",
			ImageURL:           "https://example.com/image.jpg",
			Remarks:            "This is a test aircraft",
		}
		mockAircraftArr = []model.Aircraft{
			{UserID: uint(1), AircraftModel: "Cessna 1", RegistrationNumber: "N12345", ImageURL: "https://example.com/image.jpg", Remarks: "This is a test aircraft"},
			{UserID: uint(1), AircraftModel: "Cessna 2", RegistrationNumber: "N12345", ImageURL: "https://example.com/image.jpg", Remarks: "This is a test aircraft"},
			{UserID: uint(2), AircraftModel: "Cessna 3", RegistrationNumber: "N12345", ImageURL: "https://example.com/image.jpg", Remarks: "This is a test aircraft"},
		}
	})

	AfterEach(func() {
		aircraftRepoCtrl.Finish()
		flightRepoCtrl.Finish()
	})

	Describe("InsertAircraft", func() {
		Context("when aircraft request is valid", func() {
			It("should insert aircraft and return no error", func() {
				// given
				aircraftRepoMock.EXPECT().Save(mockAircraft).Return(model.Aircraft{
					UserID:             uint(1),
					AircraftModel:      "Cessna 172",
					RegistrationNumber: "N12345",
					ImageURL:           "https://example.com/image.jpg",
					Remarks:            "This is a test aircraft",
				}, nil)

				// when
				insertedAircraft, err := aircraftService.InsertAircraft(uint(1), aircraftRequest)

				// then
				Expect(err).To(BeNil())
				Expect(insertedAircraft.UserID).To(Equal(uint(1)))
				Expect(insertedAircraft.AircraftModel).To(Equal(aircraftRequest.AircraftModel))
				Expect(insertedAircraft.RegistrationNumber).To(Equal(aircraftRequest.RegistrationNumber))
				Expect(insertedAircraft.ImageURL).To(Equal(aircraftRequest.ImageURL))
				Expect(insertedAircraft.Remarks).To(Equal(aircraftRequest.Remarks))

			})
		})
		Context("when save to database fails", func() {
			It("should return error", func() {
				// given
				aircraftRepoMock.EXPECT().Save(mockAircraft).Return(model.Aircraft{}, errors.New("failed to save aircraft"))

				// when
				insertedAircraft, err := aircraftService.InsertAircraft(uint(1), aircraftRequest)

				// then
				Expect(err.Error()).To(Equal("failed to save aircraft"))
				Expect(insertedAircraft).To(Equal(model.Aircraft{}))
			})
		})
	})

	Describe("GetUserAircraft", func() {
		Context("when user has more than one aircraft", func() {
			It("should return aircraft", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserID(uint(1)).Return(mockAircraftArr[:2], nil)

				// when
				aircraft, err := aircraftService.GetUserAircraft(uint(1))

				// then
				Expect(err).To(BeNil())
				Expect(aircraft).To(HaveLen(2))
				Expect(aircraft[0]).To(Equal(mockAircraftArr[0]))
				Expect(aircraft[1]).To(Equal(mockAircraftArr[1]))
			})
		})
		Context("when user has no aircraft", func() {
			It("should return empty slice", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserID(uint(1)).Return([]model.Aircraft{}, nil)

				// when
				aircraft, err := aircraftService.GetUserAircraft(uint(1))

				// then
				Expect(err).To(BeNil())
				Expect(len(aircraft)).To(Equal(0))
				Expect(aircraft).To(Equal([]model.Aircraft{}))
			})
		})
		Context("when user is not authorized", func() {
			It("should return error and empty slice", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserID(uint(1)).Return([]model.Aircraft{}, errors.New("failed to get aircraft"))

				// when
				aircraft, err := aircraftService.GetUserAircraft(uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to get aircraft"))
				Expect(aircraft).To(BeEmpty())
				Expect(aircraft).To(Equal([]model.Aircraft{}))
			})
		})
	})

	Describe("DeleteAircraft", func() {
		Context("when fail to count flights", func() {
			It("should return error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(0), errors.New("failed to count flights"))

				// when
				err := aircraftService.DeleteAircraft(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to count flights"))
			})
		})
		Context("when fail to delete aircraft", func() {
			It("should return error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(0), nil)
				aircraftRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(errors.New("failed to delete aircraft"))

				// when
				err := aircraftService.DeleteAircraft(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to delete aircraft"))
			})
		})
		Context("when aircraft has no flights assigned and user is authorized", func() {
			It("should return no error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(0), nil)
				aircraftRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(nil)

				// when
				err := aircraftService.DeleteAircraft(uint(1), uint(1))

				// then
				Expect(err).To(BeNil())
			})
		})
		Context("when aircraft has assigned flights", func() {
			It("should return error", func() {
				//given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(1), nil)

				// when
				err := aircraftService.DeleteAircraft(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("the plane has assigned flights"))
			})
		})
		Context("when unauthorized to delete aircraft", func() {
			It("should return error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(0), nil)
				aircraftRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(errors.New("no aircraft to delete or unauthorized to delete aircraft"))

				// when
				err := aircraftService.DeleteAircraft(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("no aircraft to delete or unauthorized to delete aircraft"))
			})
		})
	})

	Describe("UpdateAircraft", func() {
		Context("when fail to fetch aircraft", func() {
			It("should return error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(model.Aircraft{}, errors.New("failed to get aircraft"))

				// when
				updatedAircraft, err := aircraftService.UpdateAircraft(uint(1), uint(1), aircraftRequest)

				// then
				Expect(err.Error()).To(Equal("failed to get aircraft"))
				Expect(updatedAircraft).To(Equal(model.Aircraft{}))
			})
		})
		Context("when unauthorized to update contact", func() {
			It("should return error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(model.Aircraft{UserID: uint(2)}, nil)

				// when
				updatedAircraft, err := aircraftService.UpdateAircraft(uint(1), uint(1), aircraftRequest)

				// then
				Expect(err.Error()).To(Equal("unauthorized to update aircraft"))
				Expect(updatedAircraft).To(Equal(model.Aircraft{}))
			})
		})
		Context("when update fails", func() {
			It("should return error", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(mockAircraft, nil)
				aircraftRepoMock.EXPECT().Update(mockAircraft).Return(model.Aircraft{}, errors.New("failed to update aircraft"))

				// when
				updatedAircraft, err := aircraftService.UpdateAircraft(uint(1), uint(1), aircraftRequest)

				// then
				Expect(err.Error()).To(Equal("failed to update aircraft"))
				Expect(updatedAircraft).To(Equal(model.Aircraft{}))
			})
		})
		Context("when contact is updated successfully", func() {
			It("should return updated contact", func() {
				// given
				aircraftRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(mockAircraft, nil)
				aircraftRepoMock.EXPECT().Update(mockAircraft).Return(mockAircraft, nil)

				// when
				updatedAircraft, err := aircraftService.UpdateAircraft(uint(1), uint(1), aircraftRequest)

				// then
				Expect(err).To(BeNil())
				Expect(updatedAircraft).To(Equal(mockAircraft))
			})
		})
	})

	Describe("CountAircraftFlights", func() {
		Context("when fail to count flights", func() {
			It("should return error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(0), errors.New("failed to count flights"))

				// when
				count, err := aircraftService.CountAircraftFlights(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to count flights"))
				Expect(count).To(Equal(int64(0)))
			})
		})
		Context("when count flights successfully", func() {
			It("should return number of flights and no error", func() {
				// given
				flightRepoMock.EXPECT().CountByAircraftID(uint(1), uint(1)).Return(int64(5), nil)

				// when
				count, err := aircraftService.CountAircraftFlights(uint(1), uint(1))

				// then
				Expect(err).To(BeNil())
				Expect(count).To(Equal(int64(5)))
			})
		})
	})
})

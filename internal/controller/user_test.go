package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/avialog/backend/internal/utils"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("UserController", func() {
	var (
		userController   UserController
		userServiceCtrl  *gomock.Controller
		userServiceMock  *service.MockUserService
		w                *httptest.ResponseRecorder
		ctx              *gin.Context
		userRequest      dto.UserRequest
		expectedResponse dto.UserResponse
		userMock         model.User
		country          model.Country
	)

	BeforeEach(func() {
		userServiceCtrl = gomock.NewController(GinkgoT())
		userServiceMock = service.NewMockUserService(userServiceCtrl)
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
		userController = newUserController(userServiceMock)
		country = "US"

		userRequest = dto.UserRequest{
			FirstName:    utils.PtrToStr("John"),
			LastName:     utils.PtrToStr("Doe"),
			AvatarURL:    utils.PtrToStr("https://avatar.com"),
			SignatureURL: utils.PtrToStr("https://signature.com"),
			Country:      &country,
			Phone:        utils.PtrToStr("123456789"),
			Street:       utils.PtrToStr("Main St"),
			City:         utils.PtrToStr("New York"),
			Company:      utils.PtrToStr("Company"),
			Timezone:     utils.PtrToStr("UTC"),
		}
		expectedResponse = dto.UserResponse{
			FirstName:    utils.PtrToStr("John"),
			LastName:     utils.PtrToStr("Doe"),
			Email:        "test@test.com",
			AvatarURL:    utils.PtrToStr("https://avatar.com"),
			SignatureURL: utils.PtrToStr("https://signature.com"),
			Country:      &country,
			Phone:        utils.PtrToStr("123456789"),
			Street:       utils.PtrToStr("Main St"),
			City:         utils.PtrToStr("New York"),
			Company:      utils.PtrToStr("Company"),
			Timezone:     utils.PtrToStr("UTC"),
		}
		userMock = model.User{
			ID:           "1",
			FirstName:    utils.PtrToStr("John"),
			LastName:     utils.PtrToStr("Doe"),
			Email:        "test@test.com",
			AvatarURL:    utils.PtrToStr("https://avatar.com"),
			SignatureURL: utils.PtrToStr("https://signature.com"),
			Country:      &country,
			Phone:        utils.PtrToStr("123456789"),
			Street:       utils.PtrToStr("Main St"),
			City:         utils.PtrToStr("New York"),
			Company:      utils.PtrToStr("Company"),
			Timezone:     utils.PtrToStr("UTC"),
		}
	})

	AfterEach(func() {
		userServiceCtrl.Finish()
	})

	Describe("GetUser", func() {
		Context("on successful get profile", func() {
			It("should return 200 OK and user profile", func() {
				// given
				expectedJSON, err := json.Marshal(expectedResponse)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodGet, "/api/profile", nil)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req

				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				userServiceMock.EXPECT().GetUser("1").Return(userMock, nil)

				// when
				userController.GetUser(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(MatchJSON(expectedJSON))
			})
		})
		Context("on failed get profile", func() {
			It("should return 500 Internal Server Error", func() {
				// given
				req, err := http.NewRequest(http.MethodGet, "/api/profile", nil)
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				userServiceMock.EXPECT().GetUser("1").Return(model.User{}, errors.New("failed to get profile"))

				// when
				userController.GetUser(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(MatchJSON(`{"code":500,"message":"failed to get profile"}`))
			})
		})
	})
	Describe("UpdateProfile", func() {
		Context("on successful update profile", func() {
			It("should return 200 OK and updated user profile", func() {
				// given
				expectedResponseJSON, err := json.Marshal(expectedResponse)
				Expect(err).ToNot(HaveOccurred())
				userRequestJSON, err := json.Marshal(userRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/api/profile", bytes.NewBuffer(userRequestJSON))
				Expect(err).ToNot(HaveOccurred())
				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				userServiceMock.EXPECT().UpdateProfile("1", userRequest).Return(userMock, nil)
				// when
				userController.UpdateProfile(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(MatchJSON(expectedResponseJSON))
			})
		})
		Context("when binding request failed", func() {
			It("should return 400 Bad Request", func() {
				// given
				req, err := http.NewRequest(http.MethodPost, "/api/profile", bytes.NewBuffer([]byte("invalid json")))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				// when
				userController.UpdateProfile(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(MatchJSON(`{"code":400,"message":"invalid character 'i' looking for beginning of value"}`))
			})
		})
		Context("on failed update profile", func() {
			It("should return 500 Internal Server Error", func() {
				// given
				userRequestJSON, err := json.Marshal(userRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/api/profile", bytes.NewBuffer(userRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				userServiceMock.EXPECT().UpdateProfile("1", userRequest).Return(model.User{}, errors.New("failed to update profile"))

				// when
				userController.UpdateProfile(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(MatchJSON(`{"code":500,"message":"failed to update profile"}`))
			})
		})
	})
})

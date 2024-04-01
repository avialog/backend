package service

import (
	"errors"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("UserService", func() {
	var (
		userService  UserService
		userRepoCtrl *gomock.Controller
		userRepoMock *repository.MockUserRepository
		mockUser     model.User
		userRequest  dto.UserRequest
	)

	BeforeEach(func() {
		userRepoCtrl = gomock.NewController(GinkgoT())
		userRepoMock = repository.NewMockUserRepository(userRepoCtrl)
		userService = newUserService(userRepoMock, config.Config{})
		mockUser = model.User{
			ID:           "1",
			FirstName:    "test_user",
			LastName:     "test_last_name",
			Email:        "test@test.com",
			AvatarURL:    "https://example.com/avatar.jpg",
			SignatureURL: "https://example.com/signature.jpg",
			Country:      "US",
			Phone:        "1234567890",
			Street:       "1234 Main St",
			City:         "Any town",
			Company:      "Test Company",
			Timezone:     "America/New_York",
		}
		userRequest = dto.UserRequest{
			FirstName:    "test_user",
			LastName:     "test_last_name",
			AvatarURL:    "https://example.com/avatar.jpg",
			SignatureURL: "https://example.com/signature.jpg",
			Country:      "US",
			Phone:        "1234567890",
			Street:       "1234 Main St",
			City:         "Any town",
			Company:      "Test Company",
			Timezone:     "America/New_York",
		}

	})

	AfterEach(func() {
		userRepoCtrl.Finish()
	})

	Describe("GetProfile", func() {
		Context("when user exists", func() {
			It("should return user and no error", func() {
				// given
				userRepoMock.EXPECT().GetByID("1").Return(mockUser, nil)

				// when
				user, err := userService.GetUser("1")

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(mockUser))
			})
		})
		Context("when user does not exist", func() {
			It("should return error", func() {
				// given
				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, errors.New("user not found"))

				// when
				user, err := userService.GetUser("1")

				// then
				Expect(err.Error()).To(Equal("user not found"))
				Expect(user).To(Equal(model.User{}))
			})
		})
	})
	Describe("UpdateProfile", func() {
		Context("when user exists", func() {
			It("should return user and no error", func() {
				// given
				userRepoMock.EXPECT().Save(mockUser).Return(mockUser, nil)
				userRepoMock.EXPECT().GetByID("1").Return(mockUser, nil)
				// when
				user, err := userService.UpdateProfile("1", userRequest)

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(mockUser))
			})
		})
		Context("when user does not exist", func() {
			It("should return error", func() {
				// given
				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, errors.New("user not found"))
				// when
				user, err := userService.UpdateProfile("1", userRequest)

				// then
				Expect(err.Error()).To(Equal("user not found"))
				Expect(user).To(Equal(model.User{}))
			})
		})
	})
})

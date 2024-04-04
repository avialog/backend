package service

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/infrastructure"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var _ = Describe("AuthService", func() {
	var (
		authService            AuthService
		userRepoCtrl           *gomock.Controller
		userRepoMock           *repository.MockUserRepository
		authClientMock         *infrastructure.MockAuthClient
		authClientCtrl         *gomock.Controller
		userMock               model.User
		userMockWithIDandEmail model.User
	)

	BeforeEach(func() {
		authClientCtrl = gomock.NewController(GinkgoT())
		authClientMock = infrastructure.NewMockAuthClient(authClientCtrl)
		userRepoCtrl = gomock.NewController(GinkgoT())
		userRepoMock = repository.NewMockUserRepository(userRepoCtrl)
		authService = newAuthService(userRepoMock, authClientMock)
		userMockWithIDandEmail = model.User{
			ID:    "1",
			Email: "test@test.com",
		}
		userMock = model.User{
			ID:        "1",
			Email:     "test@test.com",
			FirstName: utils.String("Kate"),
		}
	})

	AfterEach(func() {
		userRepoCtrl.Finish()
		authClientCtrl.Finish()
	})

	Describe("ValidateToken", func() {
		Context("when token is valid, user is in database and didn't change email", func() {
			It("should return user", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)

				userRepoMock.EXPECT().GetByID("1").Return(userMock, nil)
				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(userMock))
			})
		})
		Context("when token is valid, user is not in database", func() {
			It("should create user and return it", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)

				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrNotFound, gorm.ErrRecordNotFound))
				userRepoMock.EXPECT().Create(userMockWithIDandEmail).Return(userMockWithIDandEmail, nil)

				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(userMockWithIDandEmail))
			})
		})
		Context("when token is valid, user is not in database, but creating user fails", func() {
			It("should return error", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)

				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrNotFound, gorm.ErrRecordNotFound))
				userRepoMock.EXPECT().Create(userMockWithIDandEmail).Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB)))
				Expect(user).To(Equal(model.User{}))
			})
		})
		Context("when token is valid, but user changed email", func() {
			It("should update email and return user", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)
				userRepoMock.EXPECT().GetByID("1").Return(model.User{ID: "1", Email: "differentEmail@test.com"}, nil)
				userRepoMock.EXPECT().Save(model.User{ID: "1", Email: "test@test.com"}).Return(model.User{ID: "1", Email: "test@test.com"}, nil)
				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(model.User{ID: "1", Email: "test@test.com"}))
			})
		})
		Context("when token is invalid", func() {
			It("should return error", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").Return(nil, errors.New("invalid token"))

				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("%w: %v", dto.ErrInternalFailure, "invalid token")))
				Expect(user).To(Equal(model.User{}))
			})
		})
		Context("when getting user from database fails and return internal error", func() {
			It("should return error", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)

				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB)))
				Expect(user).To(Equal(model.User{}))
			})
		})
		Context("when updating user email fails", func() {
			It("should return error", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)
				userRepoMock.EXPECT().GetByID("1").Return(model.User{ID: "1", Email: "differentEmail@test.com"}, nil)
				userRepoMock.EXPECT().Save(model.User{ID: "1", Email: "test@test.com"}).Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))
				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB)))
				Expect(user).To(Equal(model.User{}))
			})
		})
	})
})

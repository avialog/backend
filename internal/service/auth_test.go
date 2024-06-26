package service

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/infrastructure"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/util"
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
		userMockWithIDAndEmail model.User
		tokenExpireVerifier    infrastructure.TokenExpireVerifier
	)

	BeforeEach(func() {
		tokenExpireVerifier = func(err error) bool {
			return false
		}
		authClientCtrl = gomock.NewController(GinkgoT())
		authClientMock = infrastructure.NewMockAuthClient(authClientCtrl)
		userRepoCtrl = gomock.NewController(GinkgoT())
		userRepoMock = repository.NewMockUserRepository(userRepoCtrl)
		authService = newAuthService(userRepoMock, authClientMock, tokenExpireVerifier)
		userMockWithIDAndEmail = model.User{
			ID:    "1",
			Email: "test@test.com",
		}
		userMock = model.User{
			ID:        "1",
			Email:     "test@test.com",
			FirstName: util.String("Kate"),
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
				userRepoMock.EXPECT().Create(userMockWithIDAndEmail).Return(userMockWithIDAndEmail, nil)

				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).To(BeNil())
				Expect(user).To(Equal(userMockWithIDAndEmail))
			})
		})
		Context("when token is valid, user is not in database, but creating user fails", func() {
			It("should return error", func() {
				// given
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").
					Return(&auth.Token{Claims: map[string]interface{}{"email": "test@test.com"}, UID: "1"}, nil)

				userRepoMock.EXPECT().GetByID("1").Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrNotFound, gorm.ErrRecordNotFound))
				userRepoMock.EXPECT().Create(userMockWithIDAndEmail).Return(model.User{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

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
		Context("when token is outdated", func() {
			It("should return error", func() {
				// given
				expectedErr := errors.New("invalid token")
				authClientMock.EXPECT().VerifyIDToken(context.Background(), "test-token").Return(nil, expectedErr)
				called := false
				tokenExpireVerifier = func(err error) bool {
					Expect(err).To(Equal(expectedErr))
					called = true
					return true
				}

				authService = newAuthService(userRepoMock, authClientMock, tokenExpireVerifier)
				// when
				user, err := authService.ValidateToken(context.Background(), "test-token")

				// then
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("%w: %v", dto.ErrNotAuthorized, "invalid token")))
				Expect(user).To(Equal(model.User{}))
				Expect(called).To(BeTrue())
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

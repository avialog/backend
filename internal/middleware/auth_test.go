package middleware

import (
	"errors"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/service"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("AuthJWT Middleware", func() {
	var (
		authServiceCtrl *gomock.Controller
		authService     *service.MockAuthService
		w               *httptest.ResponseRecorder
		ctx             *gin.Context
	)

	BeforeEach(func() {
		authServiceCtrl = gomock.NewController(GinkgoT())
		authService = service.NewMockAuthService(authServiceCtrl)
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
	})

	AfterEach(func() {
		authServiceCtrl.Finish()
	})

	Context("when token is valid", func() {
		It("should pass the request to the next handler", func() {
			// given
			authService.EXPECT().ValidateToken(ctx, "valid-token").Return(model.User{ID: "1"}, nil)
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			ctx.Request.Header.Set("Authorization", "valid-token")

			// when
			AuthJWT(authService)(ctx)

			// then
			Expect(ctx.Writer.Status()).To(Equal(http.StatusOK))
			Expect(ctx.GetString("userID")).To(Equal("1"))
			Expect(ctx.IsAborted()).To(BeFalse())
			Expect(w.Code).To(Equal(http.StatusOK))

		})
	})

	Context("when token is invalid", func() {
		It("should return an error", func() {
			// given
			authService.EXPECT().ValidateToken(ctx, "invalid-token").Return(model.User{}, errors.New("invalid token"))
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			ctx.Request.Header.Set("Authorization", "invalid-token")

			// when
			AuthJWT(authService)(ctx)

			// then
			Expect(ctx.Writer.Status()).To(Equal(http.StatusInternalServerError))
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(w.Body.String()).To(Equal("{\"error\":\"invalid token\"}"))
			Expect(ctx.GetString("userID")).To(Equal(""))
			Expect(ctx.IsAborted()).To(BeTrue())
		})
	})
})

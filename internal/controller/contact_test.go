package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var _ = Describe("UserController", func() {
	var (
		contactController   ContactController
		contactServiceCtrl  *gomock.Controller
		contactServiceMock  *service.MockContactService
		w                   *httptest.ResponseRecorder
		ctx                 *gin.Context
		mockContacts        []model.Contact
		expectedContacts    []dto.ContactResponse
		contactRequest      dto.ContactRequest
		contactBeforeUpdate model.Contact
	)

	BeforeEach(func() {
		contactServiceCtrl = gomock.NewController(GinkgoT())
		contactServiceMock = service.NewMockContactService(contactServiceCtrl)
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
		contactController = newContactController(contactServiceMock)
		mockContacts = []model.Contact{
			{
				Model:        gorm.Model{ID: 1},
				UserID:       "1",
				AvatarURL:    util.String("https://test.com"),
				FirstName:    "John",
				LastName:     util.String("Doe"),
				Company:      util.String("Test Company"),
				Phone:        util.String("1234567890"),
				EmailAddress: util.String("test@test.com"),
				Note:         util.String("Test note"),
			},
			{
				Model:        gorm.Model{ID: 2},
				UserID:       "1",
				AvatarURL:    util.String("https://test.com"),
				FirstName:    "Jane",
				LastName:     util.String("Doe"),
				Company:      util.String("Test Company"),
				Phone:        util.String("1234567890"),
				EmailAddress: util.String("test2@test.com"),
				Note:         util.String("Test notes"),
			},
		}

		expectedContacts = []dto.ContactResponse{
			{
				ID:           1,
				AvatarURL:    util.String("https://test.com"),
				FirstName:    "John",
				LastName:     util.String("Doe"),
				Company:      util.String("Test Company"),
				Phone:        util.String("1234567890"),
				EmailAddress: util.String("test@test.com"),
				Note:         util.String("Test note"),
			},
			{
				ID:           2,
				AvatarURL:    util.String("https://test.com"),
				FirstName:    "Jane",
				LastName:     util.String("Doe"),
				Company:      util.String("Test Company"),
				Phone:        util.String("1234567890"),
				EmailAddress: util.String("test2@test.com"),
				Note:         util.String("Test notes"),
			},
		}
		contactRequest = dto.ContactRequest{
			AvatarURL:    util.String("https://test.com"),
			FirstName:    "John",
			LastName:     util.String("Doe"),
			Company:      util.String("Test Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("test@test.com"),
			Note:         util.String("Test note"),
		}
		contactBeforeUpdate = model.Contact{
			Model:        gorm.Model{ID: 1},
			UserID:       "5",
			AvatarURL:    util.String("https://test.com"),
			FirstName:    "John",
			LastName:     util.String("Doe"),
			Company:      util.String("Test Company"),
			Phone:        util.String("1234567890"),
			EmailAddress: util.String("test@test.com"),
			Note:         util.String("Test note"),
		}
	})

	AfterEach(func() {
		contactServiceCtrl.Finish()
	})

	Describe("GetContacts", func() {
		Context("when contacts are successfully fetched", func() {
			It("should return status 200 and contacts", func() {
				// given
				expectedContactsJSON, err := json.Marshal(expectedContacts)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodGet, "/api/contacts", nil)
				Expect(err).ToNot(HaveOccurred())
				ctx.Request = req
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().GetUserContacts("1").Return(mockContacts, nil)
				// when
				contactController.GetContacts(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(expectedContactsJSON))
			})
		})
		Context("when internal error occurs", func() {
			It("should return status 500", func() {
				// given
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().GetUserContacts("1").Return(nil, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				contactController.GetContacts(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure: invalid db"}`))
			})
		})
	})

	Describe("InsertContact", func() {
		Context("when contact is inserted successfully", func() {
			It("should return status 201 and contact", func() {
				// given
				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())
				contactResponseJSON, err := json.Marshal(expectedContacts[0])
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/api/contacts", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().InsertContact("1", contactRequest).Return(mockContacts[0], nil)
				// when
				contactController.InsertContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusCreated))
				Expect(w.Body).To(MatchJSON(contactResponseJSON))
			})
		})
		Context("when contact request missing first name", func() {
			It("should return status 400", func() {
				// given
				contactRequest.FirstName = ""
				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/api/contacts", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				// when
				contactController.InsertContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"Key: 'ContactRequest.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag"}`))
			})
		})
		Context("when internal error occurs", func() {
			It("should return status 500", func() {
				// given
				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPost, "/api/contacts", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().InsertContact("1", contactRequest).Return(model.Contact{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				contactController.InsertContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure: invalid db"}`))
			})
		})
		Context("when couldn't parse incoming request", func() {
			It("should return status 400", func() {
				// given
				req, err := http.NewRequest(http.MethodPost, "/api/contacts", bytes.NewBuffer([]byte("")))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				// when
				contactController.InsertContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"EOF"}`))
			})
		})
	})
	Describe("UpdateContact", func() {
		Context("when contact is updated successfully", func() {
			It("should return status 200 and contact", func() {
				// given
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}

				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())
				contactResponseJSON, err := json.Marshal(expectedContacts[0])
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/api/contacts/3", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().UpdateContact("1", uint(3), contactRequest).Return(contactBeforeUpdate, nil)
				// when
				contactController.UpdateContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body).To(MatchJSON(contactResponseJSON))
			})
		})
		Context("when contact request missing first name", func() {
			It("should return status 400", func() {
				// given
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}

				contactRequest.FirstName = ""

				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/api/contacts/3", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				// when
				contactController.UpdateContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"Key: 'ContactRequest.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag"}`))
			})
		})
		Context("when could not parse id", func() {
			It("should return status 400", func() {
				// given
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

				// when
				contactController.UpdateContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"strconv.ParseUint: parsing \"a\": invalid syntax"}`))
			})
		})
		Context("when contact is not found", func() {
			It("should return status 404", func() {
				// given
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}

				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/api/contacts/3", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().UpdateContact("1", uint(3), contactRequest).Return(model.Contact{}, fmt.Errorf("%w: %v", dto.ErrNotFound, gorm.ErrRecordNotFound))

				// when
				contactController.UpdateContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusNotFound))
				Expect(w.Body.String()).To(Equal(`{"code":404,"message":"not found: record not found"}`))
			})
		})
		Context("when internal error occurs", func() {
			It("should return status 500", func() {
				// given
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "3"}}

				contactRequestJSON, err := json.Marshal(contactRequest)
				Expect(err).ToNot(HaveOccurred())

				req, err := http.NewRequest(http.MethodPut, "/api/contacts/3", bytes.NewBuffer(contactRequestJSON))
				Expect(err).ToNot(HaveOccurred())

				ctx.Request = req
				ctx.Set("Content-Type", "application/json")
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				contactServiceMock.EXPECT().UpdateContact("1", uint(3), contactRequest).Return(model.Contact{}, fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				contactController.UpdateContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure: invalid db"}`))
			})
		})
	})
	Describe("DeleteContact", func() {
		Context("when contact is deleted successfully", func() {
			It("should return status 200", func() {
				// given
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
				contactServiceMock.EXPECT().DeleteContact("1", uint(1)).Return(nil)

				// when
				contactController.DeleteContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(Equal(`{"message":"Contact deleted successfully"}`))
			})
		})
		Context("when could not parse id", func() {
			It("should return status 400", func() {
				// given
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

				// when
				contactController.DeleteContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(w.Body.String()).To(Equal(`{"code":400,"message":"strconv.ParseUint: parsing \"a\": invalid syntax"}`))
			})
		})
		Context("when contact is not found", func() {
			It("should return status 404", func() {
				// given
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
				contactServiceMock.EXPECT().DeleteContact("1", uint(1)).Return(fmt.Errorf("%w: %v", dto.ErrNotFound, gorm.ErrRecordNotFound))

				// when
				contactController.DeleteContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusNotFound))
				Expect(w.Body.String()).To(Equal(`{"code":404,"message":"not found: record not found"}`))
			})
		})
		Context("when internal error occurs", func() {
			It("should return status 500", func() {
				// given
				ctx.Set("Accept", "application/json")
				ctx.Set("userID", "1")
				ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
				contactServiceMock.EXPECT().DeleteContact("1", uint(1)).Return(fmt.Errorf("%w: %v", dto.ErrInternalFailure, gorm.ErrInvalidDB))

				// when
				contactController.DeleteContact(ctx)

				// then
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				Expect(w.Body.String()).To(Equal(`{"code":500,"message":"internal failure: invalid db"}`))
			})
		})
	})
})

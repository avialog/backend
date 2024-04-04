package service

import (
	"errors"
	"github.com/avialog/backend/internal/config"
	"github.com/avialog/backend/internal/dto"
	"github.com/avialog/backend/internal/model"
	"github.com/avialog/backend/internal/repository"
	"github.com/avialog/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("ContactService", func() {
	var (
		contactService  ContactService
		contactRepoCtrl *gomock.Controller
		contactRepoMock *repository.MockContactRepository
		contactRequest  dto.ContactRequest
		mockContact     model.Contact
		mockContacts    []model.Contact
		validator       *validator.Validate
	)

	BeforeEach(func() {
		contactRepoCtrl = gomock.NewController(GinkgoT())
		contactRepoMock = repository.NewMockContactRepository(contactRepoCtrl)
		validator = utils.GetValidator()
		contactService = newContactService(contactRepoMock, config.Config{}, validator)
		contactRequest = dto.ContactRequest{
			FirstName:    "John",
			LastName:     utils.String("Doe"),
			Phone:        utils.String("1234567890"),
			AvatarURL:    utils.String("https://example.com/avatar.jpg"),
			Company:      utils.String("Example Inc"),
			EmailAddress: utils.String("test@test.com"),
			Note:         utils.String("This is a test contact"),
		}
		mockContact = model.Contact{
			UserID:       "1",
			FirstName:    "John",
			LastName:     utils.String("Doe"),
			Phone:        utils.String("1234567890"),
			AvatarURL:    utils.String("https://example.com/avatar.jpg"),
			Company:      utils.String("Example Inc"),
			EmailAddress: utils.String("test@test.com"),
			Note:         utils.String("This is a test contact"),
		}
		mockContacts = []model.Contact{
			{UserID: "1", FirstName: "John", LastName: utils.String("Doe"), Phone: utils.String("1234567890"), AvatarURL: utils.String("https://example.com/avatar.jpg"), Company: utils.String("Example Inc"), EmailAddress: utils.String("test@test.com"), Note: utils.String("This is a test contact")},
			{UserID: "1", FirstName: "Jane", LastName: utils.String("Doe"), Phone: utils.String("1234567890"), AvatarURL: utils.String("https://example.com/avatar.jpg"), Company: utils.String("Example Inc"), EmailAddress: utils.String("test@test.com"), Note: utils.String("This is a test contact")},
		}
	})

	AfterEach(func() {
		contactRepoCtrl.Finish()
	})

	Describe("InsertContact", func() {
		Context("when contact request is valid", func() {
			It("should insert contact and return no error", func() {
				// given
				contactRepoMock.EXPECT().Create(mockContact).Return(model.Contact{
					UserID:       "1",
					FirstName:    contactRequest.FirstName,
					LastName:     contactRequest.LastName,
					Phone:        contactRequest.Phone,
					AvatarURL:    contactRequest.AvatarURL,
					Company:      contactRequest.Company,
					EmailAddress: contactRequest.EmailAddress,
					Note:         contactRequest.Note,
				}, nil)

				// when
				insertedContact, err := contactService.InsertContact("1", contactRequest)

				// then
				Expect(err).To(BeNil())
				Expect(insertedContact.UserID).To(Equal("1"))
				Expect(insertedContact.FirstName).To(Equal(contactRequest.FirstName))
				Expect(insertedContact.LastName).To(Equal(contactRequest.LastName))
				Expect(insertedContact.Phone).To(Equal(contactRequest.Phone))
				Expect(insertedContact.AvatarURL).To(Equal(contactRequest.AvatarURL))
				Expect(insertedContact.Company).To(Equal(contactRequest.Company))
				Expect(insertedContact.EmailAddress).To(Equal(contactRequest.EmailAddress))
				Expect(insertedContact.Note).To(Equal(contactRequest.Note))
			})
		})
		Context("when save to database fails", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().Create(mockContact).Return(model.Contact{}, errors.New("failed to save contact"))

				// when
				insertedContact, err := contactService.InsertContact("1", contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to save contact"))
				Expect(insertedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when contact request doesn't have first name", func() {
			It("should return error", func() {
				// given
				contactRequest.FirstName = ""

				// when
				insertedContact, err := contactService.InsertContact("1", contactRequest)

				// then
				Expect(err.Error()).To(Equal("bad request: invalid data in field: FirstName"))
				Expect(insertedContact).To(Equal(model.Contact{}))
			})
		})
	})

	Describe("GetUserContacts", func() {
		Context("when user has contacts", func() {
			It("should return contacts", func() {
				// given
				contactRepoMock.EXPECT().GetByUserID("1").Return(mockContacts[:2], nil)

				// when
				contacts, err := contactService.GetUserContacts("1")

				// then
				Expect(err).To(BeNil())
				Expect(contacts).To(HaveLen(2))
				Expect(contacts[0]).To(Equal(mockContacts[0]))
				Expect(contacts[1]).To(Equal(mockContacts[1]))
			})
		})
		Context("when user has no contacts", func() {
			It("should return empty contacts", func() {
				// given
				contactRepoMock.EXPECT().GetByUserID("1").Return([]model.Contact{}, nil)

				// when
				contacts, err := contactService.GetUserContacts("1")

				// then
				Expect(err).To(BeNil())
				Expect(contacts).To(HaveLen(0))
			})
		})
		Context("when failed to get contacts", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserID("1").Return([]model.Contact{}, errors.New("failed to get contacts"))

				// when
				contacts, err := contactService.GetUserContacts("1")

				// then
				Expect(err.Error()).To(Equal("failed to get contacts"))
				Expect(contacts).To(HaveLen(0))
			})
		})
	})

	Describe("DeleteContact", func() {
		Context("when contact is deleted successfully", func() {
			It("should return no error", func() {
				// given
				contactRepoMock.EXPECT().DeleteByUserIDAndID("1", uint(1)).Return(nil)

				// when
				err := contactService.DeleteContact("1", uint(1))

				// then
				Expect(err).To(BeNil())
			})
		})
		Context("when delete operation fails", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().DeleteByUserIDAndID("1", uint(1)).Return(errors.New("failed to delete contact"))

				// when
				err := contactService.DeleteContact("1", uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to delete contact"))
			})
		})
		Context("when unauthorized to delete contact", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().DeleteByUserIDAndID("1", uint(1)).Return(errors.New("no contact to delete or unauthorized to delete contact"))

				// when
				err := contactService.DeleteContact("1", uint(1))

				// then
				Expect(err.Error()).To(Equal("no contact to delete or unauthorized to delete contact"))
			})
		})
	})
	Describe("UpdateContact", func() {
		Context("when fail to get contact", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID("1", uint(1)).Return(model.Contact{}, errors.New("failed to get contact"))

				// when
				updatedContact, err := contactService.UpdateContact("1", uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to get contact"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when update fails", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID("1", uint(1)).Return(model.Contact{UserID: "1"}, nil)
				contactRepoMock.EXPECT().Save(mockContact).Return(model.Contact{}, errors.New("failed to update contact"))

				// when
				updatedContact, err := contactService.UpdateContact("1", uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to update contact"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when contact is updated", func() {
			It("should return updated contact", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID("1", uint(1)).Return(mockContact, nil)
				contactRepoMock.EXPECT().Save(mockContact).Return(mockContact, nil)

				// when
				updatedContact, err := contactService.UpdateContact("1", uint(1), contactRequest)

				// then
				Expect(err).To(BeNil())
				Expect(updatedContact).To(Equal(mockContact))
			})
		})
		Context("when contact request doesn't have first name", func() {
			It("should return error", func() {
				// given
				contactRequest.FirstName = ""
				contactRepoMock.EXPECT().GetByUserIDAndID("1", uint(1)).Return(mockContact, nil)
				// when
				updatedContact, err := contactService.UpdateContact("1", uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("bad request: invalid data in field: FirstName"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
	})
})

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

var _ = Describe("ContactService", func() {
	var (
		contactService  ContactService
		contactRepoCtrl *gomock.Controller
		contactRepoMock *repository.MockContactRepository
		contactRequest  dto.ContactRequest
		mockContact     model.Contact
		mockContacts    []model.Contact
	)

	BeforeEach(func() {
		contactRepoCtrl = gomock.NewController(GinkgoT())
		contactRepoMock = repository.NewMockContactRepository(contactRepoCtrl)
		contactService = newContactService(contactRepoMock, dto.Config{})
		contactRequest = dto.ContactRequest{
			FirstName:    "John",
			LastName:     "Doe",
			Phone:        "1234567890",
			AvatarURL:    "https://example.com/avatar.jpg",
			Company:      "Example Inc",
			EmailAddress: "test@test.com",
			Note:         "This is a test contact",
		}
		mockContact = model.Contact{
			UserID:       uint(1),
			FirstName:    "John",
			LastName:     "Doe",
			Phone:        "1234567890",
			AvatarURL:    "https://example.com/avatar.jpg",
			Company:      "Example Inc",
			EmailAddress: "test@test.com",
			Note:         "This is a test contact",
		}
		mockContacts = []model.Contact{
			{UserID: uint(1), FirstName: "John", LastName: "Doe", Phone: "1234567890", AvatarURL: "https://example.com/avatar.jpg", Company: "Example Inc", EmailAddress: "test@test.com", Note: "This is a test contact"},
			{UserID: uint(1), FirstName: "Jane", LastName: "Doe", Phone: "1234567890", AvatarURL: "https://example.com/avatar.jpg", Company: "Example Inc", EmailAddress: "test@test.com", Note: "This is a test contact"},
		}
	})

	AfterEach(func() {
		contactRepoCtrl.Finish()
	})

	Describe("InsertContact", func() {
		Context("when contact request is valid", func() {
			It("should insert contact and return no error", func() {
				// given
				contactRepoMock.EXPECT().Save(mockContact).Return(model.Contact{
					UserID:       uint(1),
					FirstName:    contactRequest.FirstName,
					LastName:     contactRequest.LastName,
					Phone:        contactRequest.Phone,
					AvatarURL:    contactRequest.AvatarURL,
					Company:      contactRequest.Company,
					EmailAddress: contactRequest.EmailAddress,
					Note:         contactRequest.Note,
				}, nil)

				// when
				insertedContact, err := contactService.InsertContact(uint(1), contactRequest)

				// then
				Expect(err).To(BeNil())
				Expect(insertedContact.UserID).To(Equal(uint(1)))
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
				contactRepoMock.EXPECT().Save(mockContact).Return(model.Contact{}, errors.New("failed to save contact"))

				// when
				insertedContact, err := contactService.InsertContact(uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to save contact"))
				Expect(insertedContact).To(Equal(model.Contact{}))
			})
		})
	})

	Describe("GetUserContacts", func() {
		Context("when user has contacts", func() {
			It("should return contacts", func() {
				// given
				contactRepoMock.EXPECT().GetByUserID(uint(1)).Return(mockContacts[:2], nil)

				// when
				contacts, err := contactService.GetUserContacts(uint(1))

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
				contactRepoMock.EXPECT().GetByUserID(uint(1)).Return([]model.Contact{}, nil)

				// when
				contacts, err := contactService.GetUserContacts(1)

				// then
				Expect(err).To(BeNil())
				Expect(contacts).To(HaveLen(0))
			})
		})
		Context("when failed to get contacts", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserID(uint(1)).Return([]model.Contact{}, errors.New("failed to get contacts"))

				// when
				contacts, err := contactService.GetUserContacts(uint(1))

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
				contactRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(int64(1), nil)

				// when
				err := contactService.DeleteContact(uint(1), uint(1))

				// then
				Expect(err).To(BeNil())
			})
		})
		Context("when delete operation fails", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(int64(0), errors.New("failed to delete contact"))

				// when
				err := contactService.DeleteContact(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("failed to delete contact"))
			})
		})
		Context("when unauthorized to delete contact", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().DeleteByUserIDAndID(uint(1), uint(1)).Return(int64(0), nil)

				// when
				err := contactService.DeleteContact(uint(1), uint(1))

				// then
				Expect(err.Error()).To(Equal("no contact to delete or unauthorized to delete contact"))
			})
		})
	})
	Describe("UpdateContact", func() {
		Context("when fail to get contact", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(model.Contact{}, errors.New("failed to get contact"))

				// when
				updatedContact, err := contactService.UpdateContact(uint(1), uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to get contact"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when unauthorized to update contact", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(model.Contact{UserID: uint(2)}, nil)

				// when
				updatedContact, err := contactService.UpdateContact(uint(1), uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("unauthorized to update contact"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when update fails", func() {
			It("should return error", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(model.Contact{UserID: uint(1)}, nil)
				contactRepoMock.EXPECT().Update(mockContact).Return(model.Contact{}, errors.New("failed to update contact"))

				// when
				updatedContact, err := contactService.UpdateContact(uint(1), uint(1), contactRequest)

				// then
				Expect(err.Error()).To(Equal("failed to update contact"))
				Expect(updatedContact).To(Equal(model.Contact{}))
			})
		})
		Context("when contact is updated", func() {
			It("should return updated contact", func() {
				// given
				contactRepoMock.EXPECT().GetByUserIDAndID(uint(1), uint(1)).Return(mockContact, nil)
				contactRepoMock.EXPECT().Update(mockContact).Return(mockContact, nil)

				// when
				updatedContact, err := contactService.UpdateContact(uint(1), uint(1), contactRequest)

				// then
				Expect(err).To(BeNil())
				Expect(updatedContact).To(Equal(mockContact))
			})
		})
	})
})

// Code generated by MockGen. DO NOT EDIT.
// Source: contact.go
//
// Generated by this command:
//
//	mockgen -source=contact.go -destination=contact_mock.go -package service
//

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	dto "github.com/avialog/backend/internal/dto"
	model "github.com/avialog/backend/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockContactService is a mock of ContactService interface.
type MockContactService struct {
	ctrl     *gomock.Controller
	recorder *MockContactServiceMockRecorder
	isgomock struct{}
}

// MockContactServiceMockRecorder is the mock recorder for MockContactService.
type MockContactServiceMockRecorder struct {
	mock *MockContactService
}

// NewMockContactService creates a new mock instance.
func NewMockContactService(ctrl *gomock.Controller) *MockContactService {
	mock := &MockContactService{ctrl: ctrl}
	mock.recorder = &MockContactServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactService) EXPECT() *MockContactServiceMockRecorder {
	return m.recorder
}

// DeleteContact mocks base method.
func (m *MockContactService) DeleteContact(userID string, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContact", userID, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContact indicates an expected call of DeleteContact.
func (mr *MockContactServiceMockRecorder) DeleteContact(userID, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContact", reflect.TypeOf((*MockContactService)(nil).DeleteContact), userID, id)
}

// GetUserContacts mocks base method.
func (m *MockContactService) GetUserContacts(userID string) ([]model.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserContacts", userID)
	ret0, _ := ret[0].([]model.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserContacts indicates an expected call of GetUserContacts.
func (mr *MockContactServiceMockRecorder) GetUserContacts(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserContacts", reflect.TypeOf((*MockContactService)(nil).GetUserContacts), userID)
}

// InsertContact mocks base method.
func (m *MockContactService) InsertContact(userID string, contactRequest dto.ContactRequest) (model.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertContact", userID, contactRequest)
	ret0, _ := ret[0].(model.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertContact indicates an expected call of InsertContact.
func (mr *MockContactServiceMockRecorder) InsertContact(userID, contactRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertContact", reflect.TypeOf((*MockContactService)(nil).InsertContact), userID, contactRequest)
}

// UpdateContact mocks base method.
func (m *MockContactService) UpdateContact(userID string, id uint, contactRequest dto.ContactRequest) (model.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContact", userID, id, contactRequest)
	ret0, _ := ret[0].(model.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateContact indicates an expected call of UpdateContact.
func (mr *MockContactServiceMockRecorder) UpdateContact(userID, id, contactRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContact", reflect.TypeOf((*MockContactService)(nil).UpdateContact), userID, id, contactRequest)
}

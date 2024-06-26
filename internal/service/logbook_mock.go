// Code generated by MockGen. DO NOT EDIT.
// Source: logbook.go
//
// Generated by this command:
//
//	mockgen -source=logbook.go -destination=logbook_mock.go -package service
//

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"
	time "time"

	dto "github.com/avialog/backend/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockLogbookService is a mock of LogbookService interface.
type MockLogbookService struct {
	ctrl     *gomock.Controller
	recorder *MockLogbookServiceMockRecorder
}

// MockLogbookServiceMockRecorder is the mock recorder for MockLogbookService.
type MockLogbookServiceMockRecorder struct {
	mock *MockLogbookService
}

// NewMockLogbookService creates a new mock instance.
func NewMockLogbookService(ctrl *gomock.Controller) *MockLogbookService {
	mock := &MockLogbookService{ctrl: ctrl}
	mock.recorder = &MockLogbookServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogbookService) EXPECT() *MockLogbookServiceMockRecorder {
	return m.recorder
}

// DeleteLogbookEntry mocks base method.
func (m *MockLogbookService) DeleteLogbookEntry(userID string, flightID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLogbookEntry", userID, flightID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLogbookEntry indicates an expected call of DeleteLogbookEntry.
func (mr *MockLogbookServiceMockRecorder) DeleteLogbookEntry(userID, flightID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLogbookEntry", reflect.TypeOf((*MockLogbookService)(nil).DeleteLogbookEntry), userID, flightID)
}

// GetLogbookEntries mocks base method.
func (m *MockLogbookService) GetLogbookEntries(userID string, start, end time.Time) ([]dto.LogbookResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogbookEntries", userID, start, end)
	ret0, _ := ret[0].([]dto.LogbookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogbookEntries indicates an expected call of GetLogbookEntries.
func (mr *MockLogbookServiceMockRecorder) GetLogbookEntries(userID, start, end any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogbookEntries", reflect.TypeOf((*MockLogbookService)(nil).GetLogbookEntries), userID, start, end)
}

// InsertLogbookEntry mocks base method.
func (m *MockLogbookService) InsertLogbookEntry(userID string, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLogbookEntry", userID, logbookRequest)
	ret0, _ := ret[0].(dto.LogbookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLogbookEntry indicates an expected call of InsertLogbookEntry.
func (mr *MockLogbookServiceMockRecorder) InsertLogbookEntry(userID, logbookRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLogbookEntry", reflect.TypeOf((*MockLogbookService)(nil).InsertLogbookEntry), userID, logbookRequest)
}

// UpdateLogbookEntry mocks base method.
func (m *MockLogbookService) UpdateLogbookEntry(userID string, flightID uint, logbookRequest dto.LogbookRequest) (dto.LogbookResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLogbookEntry", userID, flightID, logbookRequest)
	ret0, _ := ret[0].(dto.LogbookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLogbookEntry indicates an expected call of UpdateLogbookEntry.
func (mr *MockLogbookServiceMockRecorder) UpdateLogbookEntry(userID, flightID, logbookRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLogbookEntry", reflect.TypeOf((*MockLogbookService)(nil).UpdateLogbookEntry), userID, flightID, logbookRequest)
}

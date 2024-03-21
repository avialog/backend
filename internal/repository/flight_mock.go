// Code generated by MockGen. DO NOT EDIT.
// Source: flight.go
//
// Generated by this command:
//
//	mockgen -source=flight.go -destination=flight_mock.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	model "github.com/avialog/backend/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockFlightRepository is a mock of FlightRepository interface.
type MockFlightRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFlightRepositoryMockRecorder
}

// MockFlightRepositoryMockRecorder is the mock recorder for MockFlightRepository.
type MockFlightRepositoryMockRecorder struct {
	mock *MockFlightRepository
}

// NewMockFlightRepository creates a new mock instance.
func NewMockFlightRepository(ctrl *gomock.Controller) *MockFlightRepository {
	mock := &MockFlightRepository{ctrl: ctrl}
	mock.recorder = &MockFlightRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlightRepository) EXPECT() *MockFlightRepositoryMockRecorder {
	return m.recorder
}

// CountByAircraftID mocks base method.
func (m *MockFlightRepository) CountByAircraftID(userID, aircraftID uint) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByAircraftID", userID, aircraftID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByAircraftID indicates an expected call of CountByAircraftID.
func (mr *MockFlightRepositoryMockRecorder) CountByAircraftID(userID, aircraftID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByAircraftID", reflect.TypeOf((*MockFlightRepository)(nil).CountByAircraftID), userID, aircraftID)
}

// DeleteByID mocks base method.
func (m *MockFlightRepository) DeleteByID(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockFlightRepositoryMockRecorder) DeleteByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockFlightRepository)(nil).DeleteByID), id)
}

// GetByAircraftID mocks base method.
func (m *MockFlightRepository) GetByAircraftID(aircraftID uint) ([]model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAircraftID", aircraftID)
	ret0, _ := ret[0].([]model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAircraftID indicates an expected call of GetByAircraftID.
func (mr *MockFlightRepositoryMockRecorder) GetByAircraftID(aircraftID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAircraftID", reflect.TypeOf((*MockFlightRepository)(nil).GetByAircraftID), aircraftID)
}

// GetByID mocks base method.
func (m *MockFlightRepository) GetByID(id uint) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockFlightRepositoryMockRecorder) GetByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockFlightRepository)(nil).GetByID), id)
}

// GetByUserID mocks base method.
func (m *MockFlightRepository) GetByUserID(userID uint) ([]model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", userID)
	ret0, _ := ret[0].([]model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockFlightRepositoryMockRecorder) GetByUserID(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockFlightRepository)(nil).GetByUserID), userID)
}

// Save mocks base method.
func (m *MockFlightRepository) Save(flight model.Flight) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", flight)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockFlightRepositoryMockRecorder) Save(flight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockFlightRepository)(nil).Save), flight)
}

// Update mocks base method.
func (m *MockFlightRepository) Update(flight model.Flight) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", flight)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockFlightRepositoryMockRecorder) Update(flight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFlightRepository)(nil).Update), flight)
}

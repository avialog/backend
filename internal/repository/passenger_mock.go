// Code generated by MockGen. DO NOT EDIT.
// Source: passenger.go
//
// Generated by this command:
//
//	mockgen -source=passenger.go -destination=passenger_mock.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	infrastructure "github.com/avialog/backend/internal/infrastructure"
	model "github.com/avialog/backend/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockPassengerRepository is a mock of PassengerRepository interface.
type MockPassengerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPassengerRepositoryMockRecorder
}

// MockPassengerRepositoryMockRecorder is the mock recorder for MockPassengerRepository.
type MockPassengerRepositoryMockRecorder struct {
	mock *MockPassengerRepository
}

// NewMockPassengerRepository creates a new mock instance.
func NewMockPassengerRepository(ctrl *gomock.Controller) *MockPassengerRepository {
	mock := &MockPassengerRepository{ctrl: ctrl}
	mock.recorder = &MockPassengerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPassengerRepository) EXPECT() *MockPassengerRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPassengerRepository) Create(passenger model.Passenger) (model.Passenger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", passenger)
	ret0, _ := ret[0].(model.Passenger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPassengerRepositoryMockRecorder) Create(passenger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPassengerRepository)(nil).Create), passenger)
}

// CreateTx mocks base method.
func (m *MockPassengerRepository) CreateTx(tx infrastructure.Database, passenger model.Passenger) (model.Passenger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTx", tx, passenger)
	ret0, _ := ret[0].(model.Passenger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockPassengerRepositoryMockRecorder) CreateTx(tx, passenger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockPassengerRepository)(nil).CreateTx), tx, passenger)
}

// DeleteByFlightIDTx mocks base method.
func (m *MockPassengerRepository) DeleteByFlightIDTx(tx infrastructure.Database, flightID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByFlightIDTx", tx, flightID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByFlightIDTx indicates an expected call of DeleteByFlightIDTx.
func (mr *MockPassengerRepositoryMockRecorder) DeleteByFlightIDTx(tx, flightID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByFlightIDTx", reflect.TypeOf((*MockPassengerRepository)(nil).DeleteByFlightIDTx), tx, flightID)
}

// DeleteByID mocks base method.
func (m *MockPassengerRepository) DeleteByID(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockPassengerRepositoryMockRecorder) DeleteByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockPassengerRepository)(nil).DeleteByID), id)
}

// GetByFlightID mocks base method.
func (m *MockPassengerRepository) GetByFlightID(id uint) ([]model.Passenger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFlightID", id)
	ret0, _ := ret[0].([]model.Passenger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByFlightID indicates an expected call of GetByFlightID.
func (mr *MockPassengerRepositoryMockRecorder) GetByFlightID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFlightID", reflect.TypeOf((*MockPassengerRepository)(nil).GetByFlightID), id)
}

// GetByID mocks base method.
func (m *MockPassengerRepository) GetByID(id uint) (model.Passenger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(model.Passenger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockPassengerRepositoryMockRecorder) GetByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockPassengerRepository)(nil).GetByID), id)
}

// Save mocks base method.
func (m *MockPassengerRepository) Save(passenger model.Passenger) (model.Passenger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", passenger)
	ret0, _ := ret[0].(model.Passenger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockPassengerRepositoryMockRecorder) Save(passenger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPassengerRepository)(nil).Save), passenger)
}

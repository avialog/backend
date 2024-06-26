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
	time "time"

	infrastructure "github.com/avialog/backend/internal/infrastructure"
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

// Begin mocks base method.
func (m *MockFlightRepository) Begin() infrastructure.Database {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(infrastructure.Database)
	return ret0
}

// Begin indicates an expected call of Begin.
func (mr *MockFlightRepositoryMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockFlightRepository)(nil).Begin))
}

// CountByUserIDAndAircraftID mocks base method.
func (m *MockFlightRepository) CountByUserIDAndAircraftID(userID string, aircraftID uint) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByUserIDAndAircraftID", userID, aircraftID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByUserIDAndAircraftID indicates an expected call of CountByUserIDAndAircraftID.
func (mr *MockFlightRepositoryMockRecorder) CountByUserIDAndAircraftID(userID, aircraftID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByUserIDAndAircraftID", reflect.TypeOf((*MockFlightRepository)(nil).CountByUserIDAndAircraftID), userID, aircraftID)
}

// Create mocks base method.
func (m *MockFlightRepository) Create(flight model.Flight) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", flight)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockFlightRepositoryMockRecorder) Create(flight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFlightRepository)(nil).Create), flight)
}

// CreateTx mocks base method.
func (m *MockFlightRepository) CreateTx(tx infrastructure.Database, flight model.Flight) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTx", tx, flight)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTx indicates an expected call of CreateTx.
func (mr *MockFlightRepositoryMockRecorder) CreateTx(tx, flight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTx", reflect.TypeOf((*MockFlightRepository)(nil).CreateTx), tx, flight)
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

// DeleteByIDTx mocks base method.
func (m *MockFlightRepository) DeleteByIDTx(tx infrastructure.Database, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByIDTx", tx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByIDTx indicates an expected call of DeleteByIDTx.
func (mr *MockFlightRepositoryMockRecorder) DeleteByIDTx(tx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByIDTx", reflect.TypeOf((*MockFlightRepository)(nil).DeleteByIDTx), tx, id)
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

// GetByIDTx mocks base method.
func (m *MockFlightRepository) GetByIDTx(tx infrastructure.Database, id uint) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDTx", tx, id)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDTx indicates an expected call of GetByIDTx.
func (mr *MockFlightRepositoryMockRecorder) GetByIDTx(tx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDTx", reflect.TypeOf((*MockFlightRepository)(nil).GetByIDTx), tx, id)
}

// GetByUserID mocks base method.
func (m *MockFlightRepository) GetByUserID(userID string) ([]model.Flight, error) {
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

// GetByUserIDAndDate mocks base method.
func (m *MockFlightRepository) GetByUserIDAndDate(userID string, start, end time.Time) ([]model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserIDAndDate", userID, start, end)
	ret0, _ := ret[0].([]model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserIDAndDate indicates an expected call of GetByUserIDAndDate.
func (mr *MockFlightRepositoryMockRecorder) GetByUserIDAndDate(userID, start, end any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserIDAndDate", reflect.TypeOf((*MockFlightRepository)(nil).GetByUserIDAndDate), userID, start, end)
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

// SaveTx mocks base method.
func (m *MockFlightRepository) SaveTx(tx infrastructure.Database, flight model.Flight) (model.Flight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTx", tx, flight)
	ret0, _ := ret[0].(model.Flight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveTx indicates an expected call of SaveTx.
func (mr *MockFlightRepositoryMockRecorder) SaveTx(tx, flight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTx", reflect.TypeOf((*MockFlightRepository)(nil).SaveTx), tx, flight)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: database.go
//
// Generated by this command:
//
//	mockgen -source=database.go -destination=database_mock.go -package infrastructure
//

// Package infrastructure is a generated GoMock package.
package infrastructure

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockDatabase) Commit() *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockDatabaseMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockDatabase)(nil).Commit))
}

// Create mocks base method.
func (m *MockDatabase) Create(value any) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", value)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDatabaseMockRecorder) Create(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDatabase)(nil).Create), value)
}

// Delete mocks base method.
func (m *MockDatabase) Delete(value any, where ...any) *gorm.DB {
	m.ctrl.T.Helper()
	varargs := []any{value}
	for _, a := range where {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDatabaseMockRecorder) Delete(value any, where ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{value}, where...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDatabase)(nil).Delete), varargs...)
}

// First mocks base method.
func (m *MockDatabase) First(dest any, where ...any) *gorm.DB {
	m.ctrl.T.Helper()
	varargs := []any{dest}
	for _, a := range where {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "First", varargs...)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// First indicates an expected call of First.
func (mr *MockDatabaseMockRecorder) First(dest any, where ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{dest}, where...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "First", reflect.TypeOf((*MockDatabase)(nil).First), varargs...)
}

// Rollback mocks base method.
func (m *MockDatabase) Rollback() *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockDatabaseMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockDatabase)(nil).Rollback))
}

// Save mocks base method.
func (m *MockDatabase) Save(value any) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", value)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockDatabaseMockRecorder) Save(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDatabase)(nil).Save), value)
}

// Where mocks base method.
func (m *MockDatabase) Where(query any, args ...any) *gorm.DB {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Where", varargs...)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Where indicates an expected call of Where.
func (mr *MockDatabaseMockRecorder) Where(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Where", reflect.TypeOf((*MockDatabase)(nil).Where), varargs...)
}

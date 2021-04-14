// Code generated by MockGen. DO NOT EDIT.
// Source: spooler.go

// Package mock is a generated GoMock package.
package mock

import (
	models "github.com/caos/zitadel/internal/eventstore/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockHandler is a mock of Handler interface
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// ViewModel mocks base method
func (m *MockHandler) ViewModel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ViewModel")
	ret0, _ := ret[0].(string)
	return ret0
}

// ViewModel indicates an expected call of ViewModel
func (mr *MockHandlerMockRecorder) ViewModel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ViewModel", reflect.TypeOf((*MockHandler)(nil).ViewModel))
}

// EventQuery mocks base method
func (m *MockHandler) EventQuery() (*models.SearchQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventQuery")
	ret0, _ := ret[0].(*models.SearchQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EventQuery indicates an expected call of EventQuery
func (mr *MockHandlerMockRecorder) EventQuery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventQuery", reflect.TypeOf((*MockHandler)(nil).EventQuery))
}

// Reduce mocks base method
func (m *MockHandler) Process(arg0 *models.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reduce", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reduce indicates an expected call of Reduce
func (mr *MockHandlerMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reduce", reflect.TypeOf((*MockHandler)(nil).Process), arg0)
}

// MinimumCycleDuration mocks base method
func (m *MockHandler) MinimumCycleDuration() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MinimumCycleDuration")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// MinimumCycleDuration indicates an expected call of MinimumCycleDuration
func (mr *MockHandlerMockRecorder) MinimumCycleDuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MinimumCycleDuration", reflect.TypeOf((*MockHandler)(nil).MinimumCycleDuration))
}

// MockLocker is a mock of Locker interface
type MockLocker struct {
	ctrl     *gomock.Controller
	recorder *MockLockerMockRecorder
}

// MockLockerMockRecorder is the mock recorder for MockLocker
type MockLockerMockRecorder struct {
	mock *MockLocker
}

// NewMockLocker creates a new mock instance
func NewMockLocker(ctrl *gomock.Controller) *MockLocker {
	mock := &MockLocker{ctrl: ctrl}
	mock.recorder = &MockLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLocker) EXPECT() *MockLockerMockRecorder {
	return m.recorder
}

// Renew mocks base method
func (m *MockLocker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Renew", lockerID, viewModel, waitTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// Renew indicates an expected call of Renew
func (mr *MockLockerMockRecorder) Renew(lockerID, viewModel, waitTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Renew", reflect.TypeOf((*MockLocker)(nil).Renew), lockerID, viewModel, waitTime)
}
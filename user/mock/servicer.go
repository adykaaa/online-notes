// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adykaaa/online-notes/user (interfaces: Servicer)

// Package mockuser is a generated GoMock package.
package mockuser

import (
	context "context"
	reflect "reflect"

	db "github.com/adykaaa/online-notes/db/sqlc"
	gomock "github.com/golang/mock/gomock"
)

// MockServicer is a mock of Servicer interface.
type MockServicer struct {
	ctrl     *gomock.Controller
	recorder *MockServicerMockRecorder
}

// MockServicerMockRecorder is the mock recorder for MockServicer.
type MockServicerMockRecorder struct {
	mock *MockServicer
}

// NewMockServicer creates a new mock instance.
func NewMockServicer(ctrl *gomock.Controller) *MockServicer {
	mock := &MockServicer{ctrl: ctrl}
	mock.recorder = &MockServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServicer) EXPECT() *MockServicerMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockServicer) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockServicerMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockServicer)(nil).GetUser), arg0, arg1)
}

// RegisterUser mocks base method.
func (m *MockServicer) RegisterUser(arg0 context.Context, arg1, arg2, arg3 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockServicerMockRecorder) RegisterUser(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockServicer)(nil).RegisterUser), arg0, arg1, arg2, arg3)
}

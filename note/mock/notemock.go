// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adykaaa/online-notes/server/http (interfaces: NoteService)

// Package mocknote is a generated GoMock package.
package mocknote

import (
	context "context"
	reflect "reflect"

	db "github.com/adykaaa/online-notes/db/sqlc"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockNoteService is a mock of NoteService interface.
type MockNoteService struct {
	ctrl     *gomock.Controller
	recorder *MockNoteServiceMockRecorder
}

// MockNoteServiceMockRecorder is the mock recorder for MockNoteService.
type MockNoteServiceMockRecorder struct {
	mock *MockNoteService
}

// NewMockNoteService creates a new mock instance.
func NewMockNoteService(ctrl *gomock.Controller) *MockNoteService {
	mock := &MockNoteService{ctrl: ctrl}
	mock.recorder = &MockNoteServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteService) EXPECT() *MockNoteServiceMockRecorder {
	return m.recorder
}

// CreateNote mocks base method.
func (m *MockNoteService) CreateNote(arg0 context.Context, arg1, arg2, arg3 string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNote", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNote indicates an expected call of CreateNote.
func (mr *MockNoteServiceMockRecorder) CreateNote(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNote", reflect.TypeOf((*MockNoteService)(nil).CreateNote), arg0, arg1, arg2, arg3)
}

// DeleteNote mocks base method.
func (m *MockNoteService) DeleteNote(arg0 context.Context, arg1 uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteServiceMockRecorder) DeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteService)(nil).DeleteNote), arg0, arg1)
}

// GetAllNotesFromUser mocks base method.
func (m *MockNoteService) GetAllNotesFromUser(arg0 context.Context, arg1 string) ([]db.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotesFromUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotesFromUser indicates an expected call of GetAllNotesFromUser.
func (mr *MockNoteServiceMockRecorder) GetAllNotesFromUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotesFromUser", reflect.TypeOf((*MockNoteService)(nil).GetAllNotesFromUser), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockNoteService) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockNoteServiceMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockNoteService)(nil).GetUser), arg0, arg1)
}

// RegisterUser mocks base method.
func (m *MockNoteService) RegisterUser(arg0 context.Context, arg1 *db.RegisterUserParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockNoteServiceMockRecorder) RegisterUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockNoteService)(nil).RegisterUser), arg0, arg1)
}

// UpdateNote mocks base method.
func (m *MockNoteService) UpdateNote(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 string, arg4 bool) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteServiceMockRecorder) UpdateNote(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteService)(nil).UpdateNote), arg0, arg1, arg2, arg3, arg4)
}
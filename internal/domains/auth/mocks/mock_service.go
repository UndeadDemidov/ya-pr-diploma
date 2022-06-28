// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/UndeadDemidov/ya-pr-diploma/internal/domains/auth (interfaces: CredentialManager)

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	user "github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	gomock "github.com/golang/mock/gomock"
)

// MockCredentialManager is a mock of CredentialManager interface.
type MockCredentialManager struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialManagerMockRecorder
}

// MockCredentialManagerMockRecorder is the mock recorder for MockCredentialManager.
type MockCredentialManagerMockRecorder struct {
	mock *MockCredentialManager
}

// NewMockCredentialManager creates a new mock instance.
func NewMockCredentialManager(ctrl *gomock.Controller) *MockCredentialManager {
	mock := &MockCredentialManager{ctrl: ctrl}
	mock.recorder = &MockCredentialManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialManager) EXPECT() *MockCredentialManagerMockRecorder {
	return m.recorder
}

// AddNewUser mocks base method.
func (m *MockCredentialManager) AddNewUser(arg0 context.Context, arg1 user.User, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewUser", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewUser indicates an expected call of AddNewUser.
func (mr *MockCredentialManagerMockRecorder) AddNewUser(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewUser", reflect.TypeOf((*MockCredentialManager)(nil).AddNewUser), arg0, arg1, arg2, arg3)
}

// AuthenticateUser mocks base method.
func (m *MockCredentialManager) AuthenticateUser(arg0 context.Context, arg1, arg2 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockCredentialManagerMockRecorder) AuthenticateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockCredentialManager)(nil).AuthenticateUser), arg0, arg1, arg2)
}

// GetUser mocks base method.
func (m *MockCredentialManager) GetUser(arg0 context.Context, arg1 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockCredentialManagerMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockCredentialManager)(nil).GetUser), arg0, arg1)
}
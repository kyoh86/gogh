// Code generated by MockGen. DO NOT EDIT.
// Source: command/hub.go

// Package command_test is a generated GoMock package.
package command_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v32/github"
	gogh "github.com/kyoh86/gogh/gogh"
	url "net/url"
	reflect "reflect"
)

// MockHubClient is a mock of HubClient interface
type MockHubClient struct {
	ctrl     *gomock.Controller
	recorder *MockHubClientMockRecorder
}

// MockHubClientMockRecorder is the mock recorder for MockHubClient
type MockHubClientMockRecorder struct {
	mock *MockHubClient
}

// NewMockHubClient creates a new mock instance
func NewMockHubClient(ctrl *gomock.Controller) *MockHubClient {
	mock := &MockHubClient{ctrl: ctrl}
	mock.recorder = &MockHubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHubClient) EXPECT() *MockHubClientMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockHubClient) Create(arg0 context.Context, arg1 gogh.Env, arg2 *gogh.Repo, arg3 string, arg4 *url.URL, arg5 bool) (*github.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockHubClientMockRecorder) Create(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHubClient)(nil).Create), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Fork mocks base method
func (m *MockHubClient) Fork(arg0 context.Context, arg1 gogh.Env, arg2 *gogh.Repo, arg3 string) (*gogh.Repo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fork", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*gogh.Repo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fork indicates an expected call of Fork
func (mr *MockHubClientMockRecorder) Fork(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fork", reflect.TypeOf((*MockHubClient)(nil).Fork), arg0, arg1, arg2, arg3)
}

// Repos mocks base method
func (m *MockHubClient) Repos(arg0 context.Context, arg1 gogh.Env, arg2 string, arg3, arg4, arg5, arg6 bool, arg7, arg8, arg9 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Repos", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Repos indicates an expected call of Repos
func (mr *MockHubClientMockRecorder) Repos(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Repos", reflect.TypeOf((*MockHubClient)(nil).Repos), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
}

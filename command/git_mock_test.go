// Code generated by MockGen. DO NOT EDIT.
// Source: git.go

// Package command_test is a generated GoMock package.
package command_test

import (
	gomock "github.com/golang/mock/gomock"
	url "net/url"
	reflect "reflect"
)

// MockGitClient is a mock of GitClient interface
type MockGitClient struct {
	ctrl     *gomock.Controller
	recorder *MockGitClientMockRecorder
}

// MockGitClientMockRecorder is the mock recorder for MockGitClient
type MockGitClientMockRecorder struct {
	mock *MockGitClient
}

// NewMockGitClient creates a new mock instance
func NewMockGitClient(ctrl *gomock.Controller) *MockGitClient {
	mock := &MockGitClient{ctrl: ctrl}
	mock.recorder = &MockGitClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitClient) EXPECT() *MockGitClientMockRecorder {
	return m.recorder
}

// AddRemote mocks base method
func (m *MockGitClient) AddRemote(arg0, arg1 string, arg2 *url.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRemote indicates an expected call of AddRemote
func (mr *MockGitClientMockRecorder) AddRemote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemote", reflect.TypeOf((*MockGitClient)(nil).AddRemote), arg0, arg1, arg2)
}

// Clone mocks base method
func (m *MockGitClient) Clone(arg0 string, arg1 *url.URL, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clone indicates an expected call of Clone
func (mr *MockGitClientMockRecorder) Clone(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockGitClient)(nil).Clone), arg0, arg1, arg2)
}

// Fetch mocks base method
func (m *MockGitClient) Fetch(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fetch indicates an expected call of Fetch
func (mr *MockGitClientMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockGitClient)(nil).Fetch), arg0)
}

// GetCurrentBranch mocks base method
func (m *MockGitClient) GetCurrentBranch(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentBranch", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentBranch indicates an expected call of GetCurrentBranch
func (mr *MockGitClientMockRecorder) GetCurrentBranch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentBranch", reflect.TypeOf((*MockGitClient)(nil).GetCurrentBranch), arg0)
}

// GetRemote mocks base method
func (m *MockGitClient) GetRemote(arg0, arg1 string) (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemote", arg0, arg1)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemote indicates an expected call of GetRemote
func (mr *MockGitClientMockRecorder) GetRemote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemote", reflect.TypeOf((*MockGitClient)(nil).GetRemote), arg0, arg1)
}

// GetRemotes mocks base method
func (m *MockGitClient) GetRemotes(arg0 string) (map[string]*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemotes", arg0)
	ret0, _ := ret[0].(map[string]*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemotes indicates an expected call of GetRemotes
func (mr *MockGitClientMockRecorder) GetRemotes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemotes", reflect.TypeOf((*MockGitClient)(nil).GetRemotes), arg0)
}

// Init mocks base method
func (m *MockGitClient) Init(arg0 string, arg1 bool, arg2, arg3, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockGitClientMockRecorder) Init(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockGitClient)(nil).Init), arg0, arg1, arg2, arg3, arg4)
}

// RemoveRemote mocks base method
func (m *MockGitClient) RemoveRemote(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRemote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRemote indicates an expected call of RemoveRemote
func (mr *MockGitClientMockRecorder) RemoveRemote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRemote", reflect.TypeOf((*MockGitClient)(nil).RemoveRemote), arg0, arg1)
}

// RenameRemote mocks base method
func (m *MockGitClient) RenameRemote(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameRemote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameRemote indicates an expected call of RenameRemote
func (mr *MockGitClientMockRecorder) RenameRemote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameRemote", reflect.TypeOf((*MockGitClient)(nil).RenameRemote), arg0, arg1, arg2)
}

// SetUpstreamTo mocks base method
func (m *MockGitClient) SetUpstreamTo(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUpstreamTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpstreamTo indicates an expected call of SetUpstreamTo
func (mr *MockGitClientMockRecorder) SetUpstreamTo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpstreamTo", reflect.TypeOf((*MockGitClient)(nil).SetUpstreamTo), arg0, arg1)
}

// Update mocks base method
func (m *MockGitClient) Update(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockGitClientMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGitClient)(nil).Update), arg0)
}

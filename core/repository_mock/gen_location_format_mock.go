// Code generated by MockGen. DO NOT EDIT.
// Source: repository/location_format.go
//
// Generated by this command:
//
//	mockgen -source repository/location_format.go -destination repository_mock/gen_location_format_mock.go -package repository_mock
//

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	reflect "reflect"

	repository "github.com/kyoh86/gogh/v4/core/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockLocationFormat is a mock of LocationFormat interface.
type MockLocationFormat struct {
	ctrl     *gomock.Controller
	recorder *MockLocationFormatMockRecorder
	isgomock struct{}
}

// MockLocationFormatMockRecorder is the mock recorder for MockLocationFormat.
type MockLocationFormatMockRecorder struct {
	mock *MockLocationFormat
}

// NewMockLocationFormat creates a new mock instance.
func NewMockLocationFormat(ctrl *gomock.Controller) *MockLocationFormat {
	mock := &MockLocationFormat{ctrl: ctrl}
	mock.recorder = &MockLocationFormatMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocationFormat) EXPECT() *MockLocationFormatMockRecorder {
	return m.recorder
}

// Format mocks base method.
func (m *MockLocationFormat) Format(ref repository.Location) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Format", ref)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Format indicates an expected call of Format.
func (mr *MockLocationFormatMockRecorder) Format(ref any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Format", reflect.TypeOf((*MockLocationFormat)(nil).Format), ref)
}

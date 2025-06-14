// Code generated by MockGen. DO NOT EDIT.
// Source: ./overlay/store.go
//
// Generated by this command:
//
//	mockgen -source ./overlay/store.go -destination ./overlay_mock/gen_store_mock.go -package overlay_mock
//

// Package overlay_mock is a generated GoMock package.
package overlay_mock

import (
	context "context"
	reflect "reflect"

	overlay "github.com/kyoh86/gogh/v4/core/overlay"
	gomock "go.uber.org/mock/gomock"
)

// MockOverlayStore is a mock of OverlayStore interface.
type MockOverlayStore struct {
	ctrl     *gomock.Controller
	recorder *MockOverlayStoreMockRecorder
	isgomock struct{}
}

// MockOverlayStoreMockRecorder is the mock recorder for MockOverlayStore.
type MockOverlayStoreMockRecorder struct {
	mock *MockOverlayStore
}

// NewMockOverlayStore creates a new mock instance.
func NewMockOverlayStore(ctrl *gomock.Controller) *MockOverlayStore {
	mock := &MockOverlayStore{ctrl: ctrl}
	mock.recorder = &MockOverlayStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOverlayStore) EXPECT() *MockOverlayStoreMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockOverlayStore) Load(ctx context.Context, initial func() overlay.OverlayService) (overlay.OverlayService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, initial)
	ret0, _ := ret[0].(overlay.OverlayService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockOverlayStoreMockRecorder) Load(ctx, initial any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockOverlayStore)(nil).Load), ctx, initial)
}

// Save mocks base method.
func (m *MockOverlayStore) Save(ctx context.Context, v overlay.OverlayService, force bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, v, force)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockOverlayStoreMockRecorder) Save(ctx, v, force any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOverlayStore)(nil).Save), ctx, v, force)
}

// Source mocks base method.
func (m *MockOverlayStore) Source() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Source")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Source indicates an expected call of Source.
func (mr *MockOverlayStoreMockRecorder) Source() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Source", reflect.TypeOf((*MockOverlayStore)(nil).Source))
}

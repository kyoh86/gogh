// Code generated by MockGen. DO NOT EDIT.
// Source: workspace/finder_service.go
//
// Generated by this command:
//
//	mockgen -source workspace/finder_service.go -destination workspace_mock/gen_finder_service_mock.go -package workspace_mock
//

// Package workspace_mock is a generated GoMock package.
package workspace_mock

import (
	context "context"
	iter "iter"
	reflect "reflect"

	repository "github.com/kyoh86/gogh/v4/core/repository"
	workspace "github.com/kyoh86/gogh/v4/core/workspace"
	gomock "go.uber.org/mock/gomock"
)

// MockFinderService is a mock of FinderService interface.
type MockFinderService struct {
	ctrl     *gomock.Controller
	recorder *MockFinderServiceMockRecorder
	isgomock struct{}
}

// MockFinderServiceMockRecorder is the mock recorder for MockFinderService.
type MockFinderServiceMockRecorder struct {
	mock *MockFinderService
}

// NewMockFinderService creates a new mock instance.
func NewMockFinderService(ctrl *gomock.Controller) *MockFinderService {
	mock := &MockFinderService{ctrl: ctrl}
	mock.recorder = &MockFinderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFinderService) EXPECT() *MockFinderServiceMockRecorder {
	return m.recorder
}

// FindByPath mocks base method.
func (m *MockFinderService) FindByPath(ctx context.Context, ws workspace.WorkspaceService, path string) (*repository.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPath", ctx, ws, path)
	ret0, _ := ret[0].(*repository.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPath indicates an expected call of FindByPath.
func (mr *MockFinderServiceMockRecorder) FindByPath(ctx, ws, path any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPath", reflect.TypeOf((*MockFinderService)(nil).FindByPath), ctx, ws, path)
}

// FindByReference mocks base method.
func (m *MockFinderService) FindByReference(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByReference", ctx, ws, reference)
	ret0, _ := ret[0].(*repository.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByReference indicates an expected call of FindByReference.
func (mr *MockFinderServiceMockRecorder) FindByReference(ctx, ws, reference any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByReference", reflect.TypeOf((*MockFinderService)(nil).FindByReference), ctx, ws, reference)
}

// ListAllRepository mocks base method.
func (m *MockFinderService) ListAllRepository(arg0 context.Context, arg1 workspace.WorkspaceService, arg2 workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllRepository", arg0, arg1, arg2)
	ret0, _ := ret[0].(iter.Seq2[*repository.Location, error])
	return ret0
}

// ListAllRepository indicates an expected call of ListAllRepository.
func (mr *MockFinderServiceMockRecorder) ListAllRepository(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllRepository", reflect.TypeOf((*MockFinderService)(nil).ListAllRepository), arg0, arg1, arg2)
}

// ListRepositoryInRoot mocks base method.
func (m *MockFinderService) ListRepositoryInRoot(arg0 context.Context, arg1 workspace.LayoutService, arg2 workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRepositoryInRoot", arg0, arg1, arg2)
	ret0, _ := ret[0].(iter.Seq2[*repository.Location, error])
	return ret0
}

// ListRepositoryInRoot indicates an expected call of ListRepositoryInRoot.
func (mr *MockFinderServiceMockRecorder) ListRepositoryInRoot(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRepositoryInRoot", reflect.TypeOf((*MockFinderService)(nil).ListRepositoryInRoot), arg0, arg1, arg2)
}

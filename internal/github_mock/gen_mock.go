// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/github/if.go

// Package github_mock is a generated GoMock package.
package github_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	github "github.com/kyoh86/gogh/v2/internal/github"
)

// MockAdaptor is a mock of Adaptor interface.
type MockAdaptor struct {
	ctrl     *gomock.Controller
	recorder *MockAdaptorMockRecorder
}

// MockAdaptorMockRecorder is the mock recorder for MockAdaptor.
type MockAdaptorMockRecorder struct {
	mock *MockAdaptor
}

// NewMockAdaptor creates a new mock instance.
func NewMockAdaptor(ctrl *gomock.Controller) *MockAdaptor {
	mock := &MockAdaptor{ctrl: ctrl}
	mock.recorder = &MockAdaptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdaptor) EXPECT() *MockAdaptorMockRecorder {
	return m.recorder
}

// GetHost mocks base method.
func (m *MockAdaptor) GetHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHost indicates an expected call of GetHost.
func (mr *MockAdaptorMockRecorder) GetHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHost", reflect.TypeOf((*MockAdaptor)(nil).GetHost))
}

// GetMe mocks base method.
func (m *MockAdaptor) GetMe(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockAdaptorMockRecorder) GetMe(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockAdaptor)(nil).GetMe), ctx)
}

// RepositoryCreate mocks base method.
func (m *MockAdaptor) RepositoryCreate(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryCreate", ctx, org, repo)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RepositoryCreate indicates an expected call of RepositoryCreate.
func (mr *MockAdaptorMockRecorder) RepositoryCreate(ctx, org, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryCreate", reflect.TypeOf((*MockAdaptor)(nil).RepositoryCreate), ctx, org, repo)
}

// RepositoryCreateFork mocks base method.
func (m *MockAdaptor) RepositoryCreateFork(ctx context.Context, owner, repo string, opts *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryCreateFork", ctx, owner, repo, opts)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RepositoryCreateFork indicates an expected call of RepositoryCreateFork.
func (mr *MockAdaptorMockRecorder) RepositoryCreateFork(ctx, owner, repo, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryCreateFork", reflect.TypeOf((*MockAdaptor)(nil).RepositoryCreateFork), ctx, owner, repo, opts)
}

// RepositoryCreateFromTemplate mocks base method.
func (m *MockAdaptor) RepositoryCreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *github.TemplateRepoRequest) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryCreateFromTemplate", ctx, templateOwner, templateRepo, templateRepoReq)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RepositoryCreateFromTemplate indicates an expected call of RepositoryCreateFromTemplate.
func (mr *MockAdaptorMockRecorder) RepositoryCreateFromTemplate(ctx, templateOwner, templateRepo, templateRepoReq interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryCreateFromTemplate", reflect.TypeOf((*MockAdaptor)(nil).RepositoryCreateFromTemplate), ctx, templateOwner, templateRepo, templateRepoReq)
}

// RepositoryDelete mocks base method.
func (m *MockAdaptor) RepositoryDelete(ctx context.Context, owner, repo string) (*github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryDelete", ctx, owner, repo)
	ret0, _ := ret[0].(*github.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RepositoryDelete indicates an expected call of RepositoryDelete.
func (mr *MockAdaptorMockRecorder) RepositoryDelete(ctx, owner, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryDelete", reflect.TypeOf((*MockAdaptor)(nil).RepositoryDelete), ctx, owner, repo)
}

// RepositoryGet mocks base method.
func (m *MockAdaptor) RepositoryGet(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryGet", ctx, owner, repo)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RepositoryGet indicates an expected call of RepositoryGet.
func (mr *MockAdaptorMockRecorder) RepositoryGet(ctx, owner, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryGet", reflect.TypeOf((*MockAdaptor)(nil).RepositoryGet), ctx, owner, repo)
}

// RepositoryList mocks base method.
func (m *MockAdaptor) RepositoryList(ctx context.Context, opts *github.RepositoryListOptions) ([]*github.RepositoryFragment, github.PageInfoFragment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryList", ctx, opts)
	ret0, _ := ret[0].([]*github.RepositoryFragment)
	ret1, _ := ret[1].(github.PageInfoFragment)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RepositoryList indicates an expected call of RepositoryList.
func (mr *MockAdaptorMockRecorder) RepositoryList(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryList", reflect.TypeOf((*MockAdaptor)(nil).RepositoryList), ctx, opts)
}

// UserGet mocks base method.
func (m *MockAdaptor) UserGet(ctx context.Context, user string) (*github.User, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserGet", ctx, user)
	ret0, _ := ret[0].(*github.User)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UserGet indicates an expected call of UserGet.
func (mr *MockAdaptorMockRecorder) UserGet(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserGet", reflect.TypeOf((*MockAdaptor)(nil).UserGet), ctx, user)
}

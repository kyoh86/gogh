// +build mock

// Code generated by github.com/koron/mockgo; DO NOT EDIT.

package command

import (
	"net/url"

	"github.com/github/hub/github"
	"github.com/kyoh86/gogh/gogh"
)

// Remote is a mock of remote.Remote for test.
type Remote struct {
	Repos_Ps  []*RemoteRepos_P
	Repos_Rs  []*RemoteRepos_R
	Fork_Ps   []*RemoteFork_P
	Fork_Rs   []*RemoteFork_R
	Create_Ps []*RemoteCreate_P
	Create_Rs []*RemoteCreate_R
}

// RemoteRepos_P packs input parameters of remote.Remote#Repos method.
type RemoteRepos_P struct {
	Ctx         gogh.Context
	User        string
	Own         bool
	Collaborate bool
	Member      bool
	Visibility  string
	Sort        string
	Direction   string
}

// RemoteRepos_R packs output parameters of remote.Remote#Repos method.
type RemoteRepos_R struct {
	Out0 []string
	Out1 error
}

// Repos is mock of remote.Remote#Repos method.
func (_m *Remote) Repos(ctx gogh.Context, user string, own bool, collaborate bool, member bool, visibility string, sort string, direction string) ([]string, error) {
	_m.Repos_Ps = append(_m.Repos_Ps, &RemoteRepos_P{ctx, user, own, collaborate, member, visibility, sort, direction})
	var _r *RemoteRepos_R
	_r, _m.Repos_Rs = _m.Repos_Rs[0], _m.Repos_Rs[1:]
	return _r.Out0, _r.Out1
}

// RemoteFork_P packs input parameters of remote.Remote#Fork method.
type RemoteFork_P struct {
	Ctx          gogh.Context
	Repository   *gogh.Repo
	Organization string
}

// RemoteFork_R packs output parameters of remote.Remote#Fork method.
type RemoteFork_R struct {
	Out0 *gogh.Repo
	Out1 error
}

// Fork is mock of remote.Remote#Fork method.
func (_m *Remote) Fork(ctx gogh.Context, repository *gogh.Repo, organization string) (*gogh.Repo, error) {
	_m.Fork_Ps = append(_m.Fork_Ps, &RemoteFork_P{ctx, repository, organization})
	var _r *RemoteFork_R
	_r, _m.Fork_Rs = _m.Fork_Rs[0], _m.Fork_Rs[1:]
	return _r.Out0, _r.Out1
}

// RemoteCreate_P packs input parameters of remote.Remote#Create method.
type RemoteCreate_P struct {
	Ctx         gogh.Context
	Repo        *gogh.Repo
	Description string
	Homepage    *url.URL
	Private     bool
}

// RemoteCreate_R packs output parameters of remote.Remote#Create method.
type RemoteCreate_R struct {
	Out0 *github.Repository
	Out1 error
}

// Create is mock of remote.Remote#Create method.
func (_m *Remote) Create(ctx gogh.Context, repo *gogh.Repo, description string, homepage *url.URL, private bool) (*github.Repository, error) {
	_m.Create_Ps = append(_m.Create_Ps, &RemoteCreate_P{ctx, repo, description, homepage, private})
	var _r *RemoteCreate_R
	_r, _m.Create_Rs = _m.Create_Rs[0], _m.Create_Rs[1:]
	return _r.Out0, _r.Out1
}

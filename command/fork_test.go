package command_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFork(t *testing.T) {
	ctx := context.Background()
	t.Run("CloneError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update  = false
			withSSH = false
			shallow = false
		)
		cloneErr := errors.New("clone error")

		u, _ := url.Parse("https://github.com/kyoh86/gogh")

		svc.gitClient.EXPECT().Clone(path, u, shallow).Return(cloneErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, "", repo),
			cloneErr.Error(),
		)
	})

	t.Run("UpdateError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update  = true
			withSSH = false
			shallow = false
		)
		updateErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(updateErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, "", repo),
			updateErr.Error(),
		)
	})

	t.Run("ForkError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = ""
		)
		forkErr := errors.New("fork error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(nil, forkErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			forkErr.Error(),
		)
	})

	t.Run("GetRemotesError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		getRemotesErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(nil, getRemotesErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			getRemotesErr.Error(),
		)
	})

	t.Run("RemoveRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		removeRemoteErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, gomock.Any()).Return(removeRemoteErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			removeRemoteErr.Error(),
		)
	})

	t.Run("AddRemoteError1", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		addRemoteErr := errors.New("add remote 1 error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u).Return(addRemoteErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			addRemoteErr.Error(),
		)
	})

	t.Run("AddRemoteError2", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		addRemoteErr := errors.New("add remote 2 error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(addRemoteErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			addRemoteErr.Error(),
		)
	})

	t.Run("FetchError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		fetchErr := errors.New("fetch error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(fetchErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			fetchErr.Error(),
		)
	})

	t.Run("GetBranchError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		getBranchErr := errors.New("get branch error")

		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Clone(path, u, shallow).Return(nil)
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(nil)
		svc.gitClient.EXPECT().GetCurrentBranch(path).Return("branch-test", getBranchErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			getBranchErr.Error(),
		)
	})

	t.Run("SetUpstreamError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		setUpstreamErr := errors.New("set upstream error")

		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Clone(path, u, shallow).Return(nil)
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(nil)
		svc.gitClient.EXPECT().GetCurrentBranch(path).Return("branch-test", nil)
		svc.gitClient.EXPECT().SetUpstreamTo(path, "kyoh86-tryouts/branch-test").Return(setUpstreamErr)
		assert.EqualError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			setUpstreamErr.Error(),
		)
	})

	t.Run("Clone", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)

		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Clone(path, u, shallow).Return(nil)
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(nil)
		svc.gitClient.EXPECT().GetCurrentBranch(path).Return("branch-test", nil)
		svc.gitClient.EXPECT().SetUpstreamTo(path, "kyoh86-tryouts/branch-test").Return(nil)
		assert.NoError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})

	t.Run("WithoutUpdate", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86").Return(nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "kyoh86-tryouts").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(nil)
		svc.gitClient.EXPECT().GetCurrentBranch(path).Return("branch-test", nil)
		svc.gitClient.EXPECT().SetUpstreamTo(path, "kyoh86-tryouts/branch-test").Return(nil)
		assert.NoError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})

	t.Run("WithUpdate", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		repo := mustParseRepo(t, svc.ev, "kyoh86/gogh")
		path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, svc.ev, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(ctx, svc.ev, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin": nil,
			"dummy":  nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		svc.gitClient.EXPECT().Fetch(path).Return(nil)
		svc.gitClient.EXPECT().GetCurrentBranch(path).Return("branch-test", nil)
		svc.gitClient.EXPECT().SetUpstreamTo(path, "kyoh86-tryouts/branch-test").Return(nil)
		assert.NoError(
			t,
			command.Fork(ctx, svc.ev, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})
}

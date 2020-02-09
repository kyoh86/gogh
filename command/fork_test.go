package command_test

import (
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
	t.Run("CloneError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
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
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, "", repo),
			cloneErr.Error(),
		)
	})

	t.Run("UpdateError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
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
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, "", repo),
			updateErr.Error(),
		)
	})

	t.Run("ForkError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = ""
		)
		forkErr := errors.New("fork error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(nil, forkErr)
		assert.EqualError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			forkErr.Error(),
		)
	})

	t.Run("GetRemotesError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		getRemotesErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(nil, getRemotesErr)
		assert.EqualError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			getRemotesErr.Error(),
		)
	})

	t.Run("RemoveRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		removeRemoteErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin":         nil,
			"kyoh86":         nil,
			"kyoh86-tryouts": nil,
			"dummy":          nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, gomock.Any()).Return(removeRemoteErr)
		assert.EqualError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			removeRemoteErr.Error(),
		)
	})

	t.Run("AddRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		addRemoteErr := errors.New("update error")

		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
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
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
			addRemoteErr.Error(),
		)
	})

	t.Run("Clone", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)

		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Clone(path, u, shallow).Return(nil)
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
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
		assert.NoError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})

	t.Run("WithoutUpdate", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = false
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
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
		assert.NoError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})

	t.Run("WithUpdate", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		repo := mustParseRepo(t, "kyoh86/gogh")
		path := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		const (
			update       = true
			withSSH      = false
			shallow      = false
			organization = "kyoh86-tryouts"
		)
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), os.ModePerm))

		svc.gitClient.EXPECT().Update(path).Return(nil)
		newRepo := mustParseRepo(t, "kyoh86-tryouts/gogh")
		svc.hubClient.EXPECT().Fork(svc.ctx, repo, organization).Return(newRepo, nil)
		svc.gitClient.EXPECT().GetRemotes(path).Return(map[string]*url.URL{
			"origin": nil,
			"dummy":  nil,
		}, nil)
		svc.gitClient.EXPECT().RemoveRemote(path, "origin").Return(nil)
		u1, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86", u1).Return(nil)
		u2, _ := url.Parse("https://github.com/kyoh86-tryouts/gogh")
		svc.gitClient.EXPECT().AddRemote(path, "kyoh86-tryouts", u2).Return(nil)
		assert.NoError(
			t,
			command.Fork(svc.ctx, svc.gitClient, svc.hubClient, update, withSSH, shallow, organization, repo),
		)
	})
}

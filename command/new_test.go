package command_test

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// success cases:
	// pattern: with owner (success)
	// pattern: without owner (success)
	t.Run("ProjectExists", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh", ".git")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)

		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), gogh.ErrProjectAlreadyExists.Error())
		assert.DirExists(t, local)
	})

	t.Run("LocalError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		localErr := errors.New("local error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(localErr)
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), localErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("RemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		remoteErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, remoteErr)
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), remoteErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("AddRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		addRemoteErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(addRemoteErr)
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), addRemoteErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("Success", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(nil)
		assert.NoError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		))
		assert.DirExists(t, local)
	})

	t.Run("LocalErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		localErr := errors.New("local error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(localErr) // UNDONE: Do
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), localErr.Error())
		assert.DirExists(t, local)

		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil) // UNDONE: Do
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(nil) // UNDONE: Do
		assert.NoError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		))
		assert.DirExists(t, local)
	})

	t.Run("RemoteErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		remoteErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil) // UNDONE: Do
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, remoteErr)
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), remoteErr.Error())
		assert.DirExists(t, local)

		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil) // UNDONE: Do
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(nil) // UNDONE: Do
		assert.NoError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		))
		assert.DirExists(t, local)
	})

	t.Run("AddRemoteErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		addRemoteErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil) // UNDONE: Do
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(addRemoteErr) // UNDONE: Do
		assert.EqualError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		), addRemoteErr.Error())
		assert.DirExists(t, local)

		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil) // UNDONE: Do
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(nil) // UNDONE: Do
		assert.NoError(t, command.New(
			svc.ctx,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			repo,
		))
		assert.DirExists(t, local)
	})
}

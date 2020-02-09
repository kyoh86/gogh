package command_test

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("GetRemoteError", func(t *testing.T) {
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
		localGit := filepath.Join(local, ".git")
		require.NoError(t, os.MkdirAll(localGit, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		getRemotesErr := errors.New("get remote error")
		svc.gitClient.EXPECT().GetRemotes(local).Return(nil, getRemotesErr)
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
		), getRemotesErr.Error())
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
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		initErr := errors.New("init error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(initErr)
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
		), initErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("CreateErr", func(t *testing.T) {
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
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		createErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, createErr)
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
		), createErr.Error())
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
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		addRemoteErr := errors.New("add remote error")
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
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
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
		gitClient := git.New(svc.ctx)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		localErr := errors.New("local error")
		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(func(local string) (map[string]*url.URL, error) {
			return gitClient.GetRemotes(local)
		})
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) {
			require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared))
		}).Return(localErr)
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

		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) error {
			return gitClient.Init(local, bare, template, separateGitDir, shared)
		})
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).DoAndReturn(func(local string, name string, url *url.URL) error {
			return gitClient.AddRemote(local, name, url)
		})
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

	t.Run("CreateErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.tearDown(t)
		gitClient := git.New(svc.ctx)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		createErr := errors.New("create error")
		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(func(local string) (map[string]*url.URL, error) {
			return gitClient.GetRemotes(local)
		})
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) {
			require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared))
		}).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, createErr)
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
		), createErr.Error())
		assert.DirExists(t, local)

		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) error {
			return gitClient.Init(local, bare, template, separateGitDir, shared)
		})
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).DoAndReturn(func(
			local string,
			name string,
			url *url.URL,
		) error {
			return gitClient.AddRemote(local, name, url)
		})
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
		gitClient := git.New(svc.ctx)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		repo := mustParseRepo(t, "kyoh86/gogh")
		addRemoteErr := errors.New("add remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) {
			require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared))
		}).Return(nil)
		svc.hubClient.EXPECT().Create(gomock.Any(), repo, description, homepage, private).Return(nil, nil)
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Do(func(
			local string,
			name string,
			url *url.URL,
		) {
			require.NoError(t, gitClient.AddRemote(local, name, url))
		}).Return(addRemoteErr)
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

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(func(local string) (map[string]*url.URL, error) {
			return gitClient.GetRemotes(local)
		})
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) {
			require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared))
		}).Return(nil)
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

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
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	t.Run("GetRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		localGit := filepath.Join(local, ".git")
		require.NoError(t, os.MkdirAll(localGit, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		getRemotesErr := errors.New("get remote error")
		svc.gitClient.EXPECT().GetRemotes(local).Return(nil, getRemotesErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), getRemotesErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("LocalError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		initErr := errors.New("init error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(initErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), initErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("CreateErr", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		createErr := errors.New("remote error")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(nil, createErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), createErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("AddRemoteError", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		addRemoteErr := errors.New("add remote error")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(addRemoteErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), addRemoteErr.Error())
		assert.DirExists(t, local)
	})

	t.Run("Success", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = "template"
			separateGitDir = "separeteGitDir"
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Return(nil)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Return(nil)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("LocalErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		localErr := errors.New("local error")
		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(gitClient.Init).Return(localErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), localErr.Error())
		assert.DirExists(t, local)

		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).DoAndReturn(gitClient.AddRemote)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("CreateErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		createErr := errors.New("create error")
		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(nil, createErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), createErr.Error())
		assert.DirExists(t, local)

		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).DoAndReturn(gitClient.AddRemote)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("AddRemoteErrorAndSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		addRemoteErr := errors.New("add remote error")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).Do(func(
			local string,
			name string,
			url *url.URL,
		) {
			require.NoError(t, gitClient.AddRemote(local, name, url))
		}).Return(addRemoteErr)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), addRemoteErr.Error())
		assert.DirExists(t, local)

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("SSHRemoteSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared.String()))
		require.NoError(t, gitClient.AddRemote(local, "origin", u))

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("NamedRemoteSuccess", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared.String()))
		require.NoError(t, gitClient.AddRemote(local, "kyoh86", homepage))

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).Do(func(
			directory string,
			bare bool,
			template string,
			separateGitDir string,
			shared string,
		) {
			require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared))
		}).Return(nil)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(homepage), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", homepage).Do(func(
			local string,
			name string,
			url *url.URL,
		) {
			require.NoError(t, gitClient.AddRemote(local, name, url))
		}).Return(nil)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
		assert.DirExists(t, local)
	})

	t.Run("Duplicated", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		u, _ := url.Parse("https://github.com/kyoh86/gogh")
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		svc.hubClient.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), description, homepage, private).Return(createNewRepoWithURL(u), nil)
		svc.gitClient.EXPECT().AddRemote(local, "origin", u).DoAndReturn(gitClient.AddRemote)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		svc.gitClient.EXPECT().Init(local, bare, template, separateGitDir, shared.String()).DoAndReturn(gitClient.Init)
		assert.NoError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		))
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		gitClient := new(git.Client)

		const (
			private        = false
			description    = "description"
			bare           = false
			template       = ""
			separateGitDir = ""
		)
		local := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(local, os.ModePerm))
		homepage, _ := url.Parse("https://kyoh86.dev/gogh")
		shared := command.RepoShared("false")
		spec := mustParseRepoSpec(t, "kyoh86/gogh")
		require.NoError(t, gitClient.Init(local, bare, template, separateGitDir, shared.String()))
		u, _ := url.Parse("https://github.com/kyoh86/dummy")
		require.NoError(t, gitClient.AddRemote(local, "origin", u))

		svc.gitClient.EXPECT().GetRemotes(local).DoAndReturn(gitClient.GetRemotes)
		assert.EqualError(t, command.New(
			ctx,
			svc.ev,
			svc.gitClient,
			svc.hubClient,
			private,
			description,
			homepage,
			bare,
			template,
			separateGitDir,
			shared,
			spec,
		), gogh.ErrProjectAlreadyExists.Error())
	})
}

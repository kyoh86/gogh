package command_test

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	t.Run("Pipe", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().Root().AnyTimes().Return([]string{svc.root})
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)
		svc.ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")

		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		svc.gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Pipe(svc.ctx, svc.gitClient, update, withSSH, shallow, "echo", []string{"kyoh86/gogh"}))
	})

	t.Run("Bulk", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().Root().AnyTimes().Return([]string{svc.root})
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)
		svc.ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")

		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		teardown := testutil.Stubin(t, []byte(`kyoh86/gogh`))
		defer teardown()

		svc.gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Bulk(svc.ctx, svc.gitClient, update, withSSH, shallow))
	})

	t.Run("GetAll", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().Root().AnyTimes().Return([]string{svc.root})
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)
		svc.ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")

		local1 := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		remote1, _ := url.Parse("https://github.com/kyoh86/gogh")
		local2 := filepath.Join(svc.root, "github.com", "kyoh86", "vim-gogh")
		remote2, _ := url.Parse("https://github.com/kyoh86/vim-gogh")
		shallow := false
		update := false
		withSSH := false

		svc.gitClient.EXPECT().Clone(local1, remote1, shallow).Return(nil)
		svc.gitClient.EXPECT().Clone(local2, remote2, shallow).Return(nil)
		assert.NoError(t, command.GetAll(svc.ctx, svc.gitClient, update, withSSH, shallow, []gogh.Repo{
			*mustParseRepo(t, svc.ctx, "kyoh86/gogh"),
			*mustParseRepo(t, svc.ctx, "kyoh86/vim-gogh"),
		}))
	})

	t.Run("Get", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().Root().AnyTimes().Return([]string{svc.root})
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)
		svc.ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")

		local := filepath.Join(svc.root, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		svc.gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Get(svc.ctx, svc.gitClient, update, withSSH, shallow, mustParseRepo(t, svc.ctx, "kyoh86/gogh")))
	})

	t.Run("List", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().Root().AnyTimes().Return([]string{svc.root})
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)
		svc.ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")

		proj1 := filepath.Join(svc.root, "github.com", "kyoh86", "gogh", ".git")
		require.NoError(t, os.MkdirAll(proj1, 0755))
		assert.NoError(t, command.List(svc.ctx, gogh.ShortFormatter(), false, ""))
	})
	t.Run("Root", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		svc.ctx.EXPECT().PrimaryRoot().AnyTimes().Return(svc.root)

		assert.NoError(t, command.Root(svc.ctx, false))
	})
}

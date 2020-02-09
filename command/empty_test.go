package command_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)
	ctx := &context.MockContext{
		MRoot:       []string{tmp},
		MGitHubHost: "github.com",
	}

	gitCtrl := gomock.NewController(t)

	// Assert that expected methods is invoked.
	defer gitCtrl.Finish()

	gitClient := NewMockGitClient(gitCtrl)

	t.Run("Pipe", func(t *testing.T) {
		local := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Pipe(ctx, gitClient, update, withSSH, shallow, "echo", []string{"kyoh86/gogh"}))
	})

	t.Run("Bulk", func(t *testing.T) {
		local := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		ctx.MStdin = strings.NewReader(`kyoh86/gogh`)
		gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Bulk(ctx, gitClient, update, withSSH, shallow))
	})

	t.Run("GetAll", func(t *testing.T) {
		local1 := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		remote1, _ := url.Parse("https://github.com/kyoh86/gogh")
		local2 := filepath.Join(tmp, "github.com", "kyoh86", "vim-gogh")
		remote2, _ := url.Parse("https://github.com/kyoh86/vim-gogh")
		shallow := false
		update := false
		withSSH := false

		gitClient.EXPECT().Clone(local1, remote1, shallow).Return(nil)
		gitClient.EXPECT().Clone(local2, remote2, shallow).Return(nil)
		assert.NoError(t, command.GetAll(ctx, gitClient, update, withSSH, shallow, gogh.Repos{
			*mustParseRepo(t, "kyoh86/gogh"),
			*mustParseRepo(t, "kyoh86/vim-gogh"),
		}))
	})

	t.Run("Get", func(t *testing.T) {
		local := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		remote, _ := url.Parse("https://github.com/kyoh86/gogh")
		shallow := false
		update := false
		withSSH := false

		gitClient.EXPECT().Clone(local, remote, shallow).Return(nil)
		assert.NoError(t, command.Get(ctx, gitClient, update, withSSH, shallow, mustParseRepo(t, "kyoh86/gogh")))
	})

	t.Run("List", func(t *testing.T) {
		proj1 := filepath.Join(tmp, "github.com", "kyoh86", "gogh", ".git")
		require.NoError(t, os.MkdirAll(proj1, 0755))
		assert.NoError(t, command.List(ctx, gogh.ShortFormatter(), false, false, ""))
	})
	t.Run("Root", func(t *testing.T) {
		assert.NoError(t, command.Root(ctx, false))
	})
}

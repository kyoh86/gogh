package command

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	defaultGitClient = &mockGitClient{}
	defaultHubClient = &mockHubClient{}
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := &context.MockContext{
		MRoot:       []string{tmp},
		MGitHubHost: "github.com",
		MStdin:      &bytes.Buffer{},
		MStdout:     ioutil.Discard,
		MStderr:     ioutil.Discard,
	}

	assert.NoError(t, Pipe(ctx, false, false, false, "echo", []string{"kyoh86/gogh"}))
	ctx.MStdin = strings.NewReader(`kyoh86/gogh`)
	assert.NoError(t, Bulk(ctx, false, false, false))
	mustRepo := func(name string) *gogh.Repo {
		t.Helper()
		repo, err := gogh.ParseRepo(name)
		require.NoError(t, err)
		return repo
	}
	assert.NoError(t, GetAll(ctx, false, false, false, gogh.Repos{
		*mustRepo("kyoh86/gogh"),
		*mustRepo("kyoh86/vim-gogh"),
	}))
	assert.NoError(t, Get(ctx, false, false, false, mustRepo("kyoh86/gogh")))
	assert.NoError(t, ConfigGetAll(&config.Config{}))
	assert.NoError(t, ConfigGet(&config.Config{}, "root"))
	assert.NoError(t, ConfigPut(&config.Config{}, "root", "/tmp"))
	assert.NoError(t, ConfigUnset(&config.Config{}, "root"))
	assert.NoError(t, Fork(ctx, false, false, false, false, "", "", mustRepo("kyoh86/gogh")))
	assert.NoError(t, New(ctx, false, "", &url.URL{}, false, false, false, "", "", gogh.ProjectShared("false"), mustRepo("kyoh86/gogh")))
	// assert.NoError(t, Repos(ctx, "", true, false, false, "", "", ""))
	assert.NoError(t, Setup(ctx, "gogogh", "zsh"))
	assert.NoError(t, List(ctx, gogh.ProjectListFormatShort, false, ""))
	assert.NoError(t, Where(ctx, false, false, "gogh"))
	assert.NoError(t, Root(ctx, false))
}

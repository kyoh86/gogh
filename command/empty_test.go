package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	defaultGitClient = &mockGitClient{}
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)
	ctx := &context.MockContext{
		MRoot:       []string{tmp},
		MGitHubHost: "github.com",
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
	assert.NoError(t, List(ctx, gogh.ShortFormatter(), false, false, ""))
	proj1 := filepath.Join(tmp, "github.com", "kyoh86", "gogh", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	assert.NoError(t, Root(ctx, false))
}

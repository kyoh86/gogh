package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)
	ctx := &context.MockContext{
		MRoot:       []string{tmp},
		MGitHubHost: "github.com",
	}
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
	assert.EqualError(t, GetAll(ctx, false, false, false, gogh.Repos{
		*mustRepo("https://example.com/kyoh86/gogh"),
	}), `not supported host: "example.com"`)
	assert.NoError(t, Get(ctx, false, false, false, mustRepo("kyoh86/gogh")), "success getting one")
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "github.com", "kyoh86", "gogh", ".git"), 0755))
	assert.NoError(t, Get(ctx, false, false, false, mustRepo("kyoh86/gogh")), "success getting one that is already exist")
	assert.NoError(t, Get(ctx, true, false, false, mustRepo("kyoh86/gogh")), "success updating one that is already exist")
}

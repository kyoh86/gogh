package command

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	defaultGitClient = &mockGitClient{}
	defaultHubClient = &mockHubClient{}
	root, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	defer os.RemoveAll(root)

	ctx := &context.MockContext{
		MRoot:       []string{root},
		MGitHubHost: "github.com",
	}

	mustRepo := func(name string) *gogh.Repo {
		t.Helper()
		repo, err := gogh.ParseRepo(name)
		require.NoError(t, err)
		return repo
	}
	assert.NoError(t, New(
		ctx,
		false,
		"",
		&url.URL{},
		false,
		false,
		false,
		"",
		"",
		gogh.ProjectShared("false"),
		mustRepo("kyoh86/gogh"),
	))

	assert.EqualError(t, New(
		ctx,
		false,
		"",
		&url.URL{},
		false,
		false,
		false,
		"",
		"",
		gogh.ProjectShared("false"),
		mustRepo("kyoh86/gogh"),
	), "project already exists")
}

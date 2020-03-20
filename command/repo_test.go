package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/require"
)

func mustParseRepoSpec(t *testing.T, name string) *gogh.RepoSpec {
	t.Helper()
	var spec gogh.RepoSpec
	require.NoError(t, spec.Set(name))
	return &spec
}

func mustParseRepo(t *testing.T, ev gogh.Env, name string) *gogh.Repo {
	t.Helper()
	repo, err := gogh.ParseRepo(ev, name)
	require.NoError(t, err)
	return repo
}

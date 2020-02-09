package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/require"
)

func mustParseRepo(t *testing.T, name string) *gogh.Repo {
	t.Helper()
	repo, err := gogh.ParseRepo(name)
	require.NoError(t, err)
	return repo
}

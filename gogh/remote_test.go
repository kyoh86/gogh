package gogh

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseURL(t *testing.T, urlString string) *url.URL {
	t.Helper()
	u, err := url.Parse(urlString)
	require.NoError(t, err)
	return u
}

func TestNewRepositoryGitHub(t *testing.T) {
	var (
		repo RemoteRepo
		err  error
	)
	ctx := &implContext{}

	repo, err = NewRepository(ctx, parseURL(t, "https://github.com/motemen/pusheen-explorer"))
	require.NoError(t, err)
	assert.True(t, repo.IsValid())

	repo, err = NewRepository(ctx, parseURL(t, "https://github.com/motemen/pusheen-explorer/"))
	require.NoError(t, err)
	assert.True(t, repo.IsValid())

	repo, err = NewRepository(ctx, parseURL(t, "https://github.com/motemen/pusheen-explorer/blob/master/README.md"))
	require.NoError(t, err)
	assert.False(t, repo.IsValid())
}

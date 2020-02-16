package gogh_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseRepo(t *testing.T, ctx gogh.Context, name string) *gogh.Repo {
	t.Helper()
	repo, err := gogh.ParseRepo(ctx, name)
	require.NoError(t, err)
	return repo
}

func TestRepoParse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := NewMockContext(ctrl)
	ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")
	ctx.EXPECT().GitHubUser().AnyTimes().Return("kyoh86")

	t.Run("full HTTPS URL", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "https://github.com/kyoh86/pusheen-explorer")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", repo.String())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", repo.URL(false).String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.URL(true).String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "git@github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "git@github.com:/kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("owner/name repo", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "kyoh86/gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/gogh", repo.String())
	})

	t.Run("name only repo", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ctx, "gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/gogh", repo.String())
	})

	t.Run("fail when invalid url given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "://////")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when empty owner is given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when empty name is given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "test/")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name contains invalid character", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "kyoh_86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name starts with hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "-kyoh86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name ends with hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "kyoh86-/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when project name contains invalid character", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "kyoh86/foo,bar")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name contains double hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "kyoh--86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when url has no path", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "https://github.com/")
		assert.EqualError(t, err, "empty project name")
		assert.Nil(t, r)
	})

	t.Run("fail when url has subfolder", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, "https://github.com/kyoh86/gogh/blob/master/gogh/repo.go")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail to parse `dot`", func(t *testing.T) {
		r, err := gogh.ParseRepo(ctx, ".")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})
}

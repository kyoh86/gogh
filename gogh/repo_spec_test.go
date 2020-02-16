package gogh_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := NewMockContext(ctrl)
	ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")
	ctx.EXPECT().GitHubUser().AnyTimes().Return("kyoh86")

	t.Run("full HTTPS URL", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("https://github.com/kyoh86/pusheen-explorer"))
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("git@github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "git@github.com:kyoh86/pusheen-explorer.git", spec.String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("git@github.com:/kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "git@github.com:/kyoh86/pusheen-explorer.git", spec.String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "github.com:kyoh86/pusheen-explorer.git", spec.String())
	})

	t.Run("owner/name spec", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("kyoh86/gogh"))
		assert.Equal(t, "kyoh86/gogh", spec.String())
	})

	t.Run("name only spec", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		require.NoError(t, spec.Set("gogh"))
		assert.Equal(t, "gogh", spec.String())
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

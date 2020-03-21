package gogh_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseRepoSpec(t *testing.T, name string) *gogh.RepoSpec {
	t.Helper()
	var spec gogh.RepoSpec
	require.NoError(t, spec.Set(name))
	return &spec
}

func TestRepoParse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ev := NewMockEnv(ctrl)
	ev.EXPECT().GithubHost().AnyTimes().Return("github.com")

	t.Run("full HTTPS URL", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ev, "https://github.com/kyoh86/pusheen-explorer")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", repo.String())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", repo.URL(false).String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.URL(true).String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ev, "git@github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ev, "git@github.com:/kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ev, "github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/pusheen-explorer", repo.FullName())
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer", repo.String())
	})

	t.Run("owner/name repo", func(t *testing.T) {
		repo, err := gogh.ParseRepo(ev, "kyoh86/gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/gogh", repo.String())
	})

	t.Run("name only repo", func(t *testing.T) {
		ev.EXPECT().GithubUser().Return("kyoh86")
		repo, err := gogh.ParseRepo(ev, "gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", repo.FullName())
		assert.Equal(t, "https://github.com/kyoh86/gogh", repo.String())
	})

	t.Run("fail when invalid url given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "://////")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when empty owner is given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when empty name is given", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "test/")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name contains invalid character", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "kyoh_86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name starts with hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "-kyoh86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name ends with hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "kyoh86-/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when project name contains invalid character", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "kyoh86/foo,bar")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when owner name contains double hyphen", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "kyoh--86/test")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail when url has no path", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "https://github.com/")
		assert.EqualError(t, err, "project name is empty")
		assert.Nil(t, r)
	})

	t.Run("fail when url has subfolder", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, "https://github.com/kyoh86/gogh/blob/master/gogh/repo.go")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("fail to parse `dot`", func(t *testing.T) {
		r, err := gogh.ParseRepo(ev, ".")
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})
}

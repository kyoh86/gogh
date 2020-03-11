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
	env := NewMockEnv(ctrl)
	env.EXPECT().GithubHost().AnyTimes().Return("github.com")

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
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("://////"))
	})

	t.Run("fail when empty owner is given", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("/test"))
	})

	t.Run("fail when empty name is given", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("test/"))
	})

	t.Run("fail when owner name contains invalid character", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("kyoh_86/test"))
	})

	t.Run("fail when owner name starts with hyphen", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("-kyoh86/test"))
	})

	t.Run("fail when owner name ends with hyphen", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("kyoh86-/test"))
	})

	t.Run("fail when project name contains invalid character", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("kyoh86/foo,bar"))
	})

	t.Run("fail when owner name contains double hyphen", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("kyoh--86/test"))
	})

	t.Run("fail when url has no path", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.EqualError(t, spec.Set("https://github.com/"), "empty project name")
	})

	t.Run("fail when url has subfolder", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("https://github.com/kyoh86/gogh/blob/master/gogh/repo.go"))
	})

	t.Run("fail to parse `dot`", func(t *testing.T) {
		spec := new(gogh.RepoSpec)
		assert.Error(t, spec.Set("."))
	})
}

func TestRepoSpecs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	env := NewMockEnv(ctrl)
	env.EXPECT().GithubHost().AnyTimes().Return("github.com")

	var specs gogh.RepoSpecs
	require.NoError(t, specs.Set("https://github.com/kyoh86/pusheen-explorer"), "full HTTPS URL")
	require.Len(t, specs, 1)
	assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", specs[0].String())

	require.NoError(t, specs.Set("git@github.com:kyoh86/pusheen-explorer.git"), "scp like URL 1")
	require.Len(t, specs, 2)
	assert.Equal(t, "git@github.com:kyoh86/pusheen-explorer.git", specs[1].String())

	require.NoError(t, specs.Set("gogh"))
	require.Len(t, specs, 3)
	assert.Equal(t, "gogh", specs[2].String())

	require.Error(t, specs.Set("://////"), "fail when invalid url given")
	require.Len(t, specs, 3)

	assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer,git@github.com:kyoh86/pusheen-explorer.git,gogh", specs.String())

	repos, err := specs.Validate(env)
	require.NoError(t, err)
	require.Len(t, repos, 3)
	assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", repos[0].String())
	assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", repos[1].String())
	assert.Equal(t, "https://github.com/kyoh86/gogh", repos[2].String())
}

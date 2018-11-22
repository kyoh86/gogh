package gogh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoName(t *testing.T) {
}

func TestRepoSpec(t *testing.T) {
	t.Run("full HTTPS URL", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := NewSpec("https://github.com/kyoh86/pusheen-explorer")
		require.NoError(t, err)
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.URL(ctx, false).String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", spec.URL(ctx, true).String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := NewSpec("git@github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "git@github.com:kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := NewSpec("git@github.com:/kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "git@github.com:/kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := NewSpec("github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "github.com:kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("owner/name spec", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := NewSpec("kyoh86/gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})

	t.Run("name only spec", func(t *testing.T) {
		ctx := &implContext{userName: "kyoh86"}
		spec, err := NewSpec("gogh")
		require.NoError(t, err)
		assert.Equal(t, "gogh", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})

	t.Run("fail when invalid url given", func(t *testing.T) {
		_, err := NewSpec("://////")
		assert.NotNil(t, err)
	})
}

func TestRepoSpecs(t *testing.T) {
	var specs RepoSpecs
	require.True(t, specs.IsCumulative())
	assert.Empty(t, specs.String())
	assert.NoError(t, specs.Set("kyoh86/gogh"), "owner/name spec")
	assert.NoError(t, specs.Set("gogh"), "name only spec")
	assert.NotNil(t, specs.Set("://////"), "fail when invalid url given")
	assert.Len(t, specs, 2)
	assert.Equal(t, "kyoh86/gogh", specs[0].String())
	assert.Equal(t, "gogh", specs[1].String())
	assert.Equal(t, "kyoh86/gogh,gogh", specs.String())
}

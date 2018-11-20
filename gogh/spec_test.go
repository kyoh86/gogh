package gogh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec(t *testing.T) {
	t.Run("full HTTPS URL", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &implContext{}
		require.NoError(t, spec.Set("https://github.com/kyoh86/pusheen-explorer"))
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.URL(ctx, false).String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", spec.URL(ctx, true).String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &implContext{}
		require.NoError(t, spec.Set("git@github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &implContext{}
		require.NoError(t, spec.Set("git@github.com:/kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &implContext{}
		require.NoError(t, spec.Set("github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("owner/name spec", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &implContext{}
		require.NoError(t, spec.Set("kyoh86/gogh"))
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})

	t.Run("name only spec", func(t *testing.T) {
		spec := new(RepoSpec)
		require.NoError(t, spec.Set("gogh"))
		ctx := &implContext{userName: "kyoh86"}
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})
}

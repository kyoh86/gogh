package gogh

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec(t *testing.T) {
	noErrorURLString := func(t *testing.T) func(u *url.URL, err error) string {
		return func(u *url.URL, err error) string {
			t.Helper()
			require.NoError(t, err)
			return u.String()
		}
	}
	t.Run("full HTTPS URL", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &mockContext{}
		require.NoError(t, spec.Set("https://github.com/kyoh86/pusheen-explorer"))
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", noErrorURLString(t)(spec.URL(ctx, false)))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", noErrorURLString(t)(spec.URL(ctx, true)))
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &mockContext{}
		require.NoError(t, spec.Set("git@github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", noErrorURLString(t)(spec.URL(ctx, false)))
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &mockContext{}
		require.NoError(t, spec.Set("git@github.com:/kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", noErrorURLString(t)(spec.URL(ctx, false)))
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &mockContext{}
		require.NoError(t, spec.Set("github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer.git", noErrorURLString(t)(spec.URL(ctx, false)))
	})

	t.Run("owner/name spec", func(t *testing.T) {
		spec := new(RepoSpec)
		ctx := &mockContext{}
		require.NoError(t, spec.Set("kyoh86/gogh"))
		assert.Equal(t, "https://github.com/kyoh86/gogh", noErrorURLString(t)(spec.URL(ctx, false)))
	})

	t.Run("name only spec", func(t *testing.T) {
		spec := new(RepoSpec)
		require.NoError(t, spec.Set("gogh"))
		ctx := &mockContext{userName: "kyoh86"}
		assert.Equal(t, "https://github.com/kyoh86/gogh", noErrorURLString(t)(spec.URL(ctx, false)))
	})
}

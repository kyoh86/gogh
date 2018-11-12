package repo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec(t *testing.T) {
	t.Run("full HTTPS URL", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("https://github.com/kyoh86/pusheen-explorer"))
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.URL().String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", spec.SSHURL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("git@github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("git@github.com:/kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("github.com:kyoh86/pusheen-explorer.git"))
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer.git", spec.URL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})

	t.Run("owner/name spec", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("kyoh86/gogh"))
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})

	os.Setenv("GITHUB_USER", "kyoh86")
	t.Run("name only spec", func(t *testing.T) {
		spec := new(Spec)
		require.NoError(t, spec.Set("gogh"))
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL().String())
		assert.Equal(t, "github.com", spec.URL().Host)
	})
}

package gogh

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	resetEnv := func(t *testing.T) {
		t.Helper()
		for _, key := range envNames {
			require.NoError(t, os.Setenv(key, ""))
		}
	}

	t.Run("get context without roots", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))

		cfg, err := GetEnvarConfig()
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", cfg.GitHubToken())
		assert.Equal(t, "hostx1", cfg.GitHubHost())
		assert.Equal(t, "kyoh86", cfg.GitHubUser())
		assert.Equal(t, "trace", cfg.LogLevel())
	})

	t.Run("expect to get invalid user name", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubUser, "-kyoh88"))
		cfg, err := GetEnvarConfig()
		require.NoError(t, err)
		require.NotNil(t, ValidateContext(cfg))
	})

	t.Run("get root paths", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))
		cfg, err := GetEnvarConfig()
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
	})

	t.Run("expects roots are not duplicated", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar:/bar:/foo"))
		cfg, err := GetEnvarConfig()
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
	})
}

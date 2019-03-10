package gogh

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
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
		require.NoError(t, os.Setenv(envLogLevel, "trace"))
		require.NoError(t, os.Setenv(envGHEHosts, "example.com example.com:9999"))

		ctx, err := CurrentContext(nil)
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", ctx.GitHubToken())
		assert.Equal(t, "hostx1", ctx.GitHubHost())
		assert.Equal(t, "kyoh86", ctx.UserName())
		assert.Equal(t, "trace", ctx.LogLevel())
		assert.Equal(t, []string{"example.com", "example.com:9999"}, ctx.GHEHosts())
	})

	t.Run("get GitHub token", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGitHubToken, "tokenx2"))
		assert.Equal(t, "tokenx2", getGitHubToken())
	})

	t.Run("get GitHub host", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGitHubHost, "hostx2"))
		assert.Equal(t, "hostx2", getGitHubHost())
	})

	t.Run("get GitHub user name", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGitHubUser, "kyoh87"))
		assert.Equal(t, "kyoh87", getUserName())
	})

	t.Run("get OS user name", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envUserName, "kyoh88"))
		assert.Equal(t, "kyoh88", getUserName())
	})

	t.Run("expect to fail to get user name from anywhere", func(t *testing.T) {
		resetEnv(t)
		assert.Panics(t, func() { getUserName() })
	})

	t.Run("expect to fail to get log level from anywhere", func(t *testing.T) {
		resetEnv(t)
		assert.Equal(t, "warn", getLogLevel())
	})

	t.Run("get root paths", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envRoot, "/foo:/bar"))
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, rts)
	})

	t.Run("get root paths from GOPATH", func(t *testing.T) {
		resetEnv(t)
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{filepath.Join(build.Default.GOPATH, "src")}, rts)
	})

	t.Run("expects roots are not duplicated", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envRoot, "/foo:/bar:/bar:/foo"))
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, rts)
	})

	t.Run("expects GHE hosts are not duplicated", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGHEHosts, "example.com example.com:9999 example.com:9999 example.com"))
		assert.Equal(t, []string{"example.com", "example.com:9999"}, getGHEHosts())
	})
}

package gogh

import (
	"context"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetEnv(t *testing.T) {
	t.Helper()
	for _, key := range envNames {
		require.NoError(t, os.Setenv(key, ""))
	}
}

func TestContext(t *testing.T) {
	t.Run("get context from envar (without configuration)", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGitHubToken, "tokenx2")) // alt: never used
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGitHubHost, "hostx2")) // alt: never used
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envGitHubUser, "kyoh87")) // alt: never used
		require.NoError(t, os.Setenv(envUserName, "kyoh88"))   // alt: never used
		require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))
		ctx, err := CurrentContext(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", ctx.GitHubToken())
		assert.Equal(t, "hostx1", ctx.GitHubHost())
		assert.Equal(t, "kyoh86", ctx.GitHubUser())
		assert.Equal(t, "trace", ctx.LogLevel())
		assert.Equal(t, []string{"/foo", "/bar"}, ctx.Roots())
	})

	t.Run("get context from envar (with configuration)", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGitHubToken, "tokenx2")) // alt: never used
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGitHubHost, "hostx2")) // alt: never used
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envGitHubUser, "kyoh87")) // alt: never used
		require.NoError(t, os.Setenv(envUserName, "kyoh88"))   // alt: never used
		require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))
		ctx, err := CurrentContext(context.Background(), Config{ // config: never used
			confKeyGitHubToken: "tokenx2",
			confKeyGitHubHost:  "hostx2",
			confKeyGitHubUser:  "kyoh87",
			confKeyLogLevel:    "error",
			confKeyRoot:        []string{"/baz", "/bux"},
		})
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", ctx.GitHubToken())
		assert.Equal(t, "hostx1", ctx.GitHubHost())
		assert.Equal(t, "kyoh86", ctx.GitHubUser())
		assert.Equal(t, "trace", ctx.LogLevel())
		assert.Equal(t, []string{"/foo", "/bar"}, ctx.Roots())
	})

	t.Run("get context from configuration", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGitHubToken, "tokenx2")) // alt: never used
		require.NoError(t, os.Setenv(envGitHubHost, "hostx2"))   // alt: never used
		require.NoError(t, os.Setenv(envGitHubUser, "kyoh87"))   // alt: never used
		require.NoError(t, os.Setenv(envUserName, "kyoh88"))     // alt: never used
		ctx, err := CurrentContext(context.Background(), Config{
			confKeyGitHubToken: "tokenx1",
			confKeyGitHubHost:  "hostx1",
			confKeyGitHubUser:  "kyoh86",
			confKeyLogLevel:    "error",
			confKeyRoot:        []string{"/baz", "/bux"},
		})
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", ctx.GitHubToken())
		assert.Equal(t, "hostx1", ctx.GitHubHost())
		assert.Equal(t, "kyoh86", ctx.GitHubUser())
		assert.Equal(t, "error", ctx.LogLevel())
		assert.Equal(t, []string{"/baz", "/bux"}, ctx.Roots())
	})

	t.Run("get context from alt envars", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envUserName, "kyoh87")) // low priority: never used
		ctx, err := CurrentContext(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", ctx.GitHubToken())
		assert.Equal(t, "hostx1", ctx.GitHubHost())
		assert.Equal(t, "kyoh86", ctx.GitHubUser())
	})

	t.Run("get default context", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		ctx, err := CurrentContext(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, "", ctx.GitHubToken())
		assert.Equal(t, DefaultHost, ctx.GitHubHost())
		assert.Equal(t, "warn", ctx.LogLevel())
		assert.Equal(t, []string{filepath.Join(build.Default.GOPATH, "src")}, ctx.Roots())
	})

	t.Run("expect to fail to get user name from anywhere", func(t *testing.T) {
		resetEnv(t)
		_, err := CurrentContext(context.Background(), nil)
		require.EqualError(t, err, "failed to find user name. set GOGH_GITHUB_USER in environment variable")
	})

	t.Run("expect to get invalid user name", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envUserName, "-kyoh88"))
		_, err := CurrentContext(context.Background(), nil)
		require.EqualError(t, err, "owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen")
	})

	t.Run("expects roots are not duplicated", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar:/bar:/foo"))
		assert.Equal(t, []string{"/foo", "/bar"}, getRoots(nil))
	})
}

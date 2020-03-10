package config

import (
	"bytes"
	"os"
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

func TestDefaultConfig(t *testing.T) {
	resetEnv(t)
	cfg := DefaultConfig()
	assert.Equal(t, "", cfg.GithubToken())
	assert.Equal(t, "github.com", cfg.GithubHost())
	assert.Equal(t, "", cfg.GithubUser())
	assert.NotEmpty(t, cfg.Root())
	assert.NotEmpty(t, cfg.PrimaryRoot())
}

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)
		cfg, err := LoadConfig(bytes.NewBufferString(`
root:
- /foo
- /bar


github:
  token: tokenx1
  user: kyoh86
  host: hostx1
`))
		require.NoError(t, err)
		assert.Equal(t, "", cfg.GithubToken(), "token should not be saved in file")
		assert.Equal(t, "hostx1", cfg.GithubHost())
		assert.Equal(t, "kyoh86", cfg.GithubUser())
		assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
		assert.Equal(t, "/foo", cfg.PrimaryRoot())
	})
	t.Run("invalid format", func(t *testing.T) {
		resetEnv(t)
		_, err := LoadConfig(bytes.NewBufferString(`{`))
		assert.Error(t, err)
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)

		var buf bytes.Buffer
		var cfg Config
		cfg.GitHub.Token = "token1"
		cfg.GitHub.Host = "hostx1"
		cfg.GitHub.User = "kyoh86"
		cfg.VRoot = []string{"/foo", "/bar"}

		require.NoError(t, SaveConfig(&buf, &cfg))

		output := buf.String()
		assert.Contains(t, output, "root:")
		assert.Contains(t, output, "- /foo")
		assert.Contains(t, output, "- /bar")
		assert.Contains(t, output, "github:")
		assert.NotContains(t, output, "tokenx1")
		assert.Contains(t, output, "  user: kyoh86")
		assert.Contains(t, output, "  host: hostx1")
	})
}

func TestGetEnvarConfig(t *testing.T) {
	resetEnv(t)
	require.NoError(t, os.Setenv(envGoghGithubToken, "tokenx1"))
	require.NoError(t, os.Setenv(envGoghGithubHost, "hostx1"))
	require.NoError(t, os.Setenv(envGoghGithubUser, "kyoh86"))
	require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar:/bar:/foo"))
	cfg, err := GetEnvarConfig()
	require.NoError(t, err)
	assert.Equal(t, "tokenx1", cfg.GithubToken())
	assert.Equal(t, "hostx1", cfg.GithubHost())
	assert.Equal(t, "kyoh86", cfg.GithubUser())
	assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root(), "expects roots are not duplicated")
	assert.Equal(t, "/foo", cfg.PrimaryRoot())
}

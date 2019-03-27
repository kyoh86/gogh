package config

import (
	"bytes"
	"log"
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

	t.Run("merging priority", func(t *testing.T) {
		resetEnv(t)

		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))
		require.NoError(t, os.Setenv(envGoghLogDate, "1"))
		require.NoError(t, os.Setenv(envGoghLogTime, "1"))
		require.NoError(t, os.Setenv(envGoghLogMicroSeconds, "1"))
		require.NoError(t, os.Setenv(envGoghLogLongFile, "1"))
		require.NoError(t, os.Setenv(envGoghLogShortFile, "1"))
		require.NoError(t, os.Setenv(envGoghLogUTC, "1"))
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))

		cfg1, err := GetEnvarConfig()
		require.NoError(t, err)

		t.Run("full overwritten config", func(t *testing.T) {
			require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx2"))
			require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx2"))
			require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh87"))
			require.NoError(t, os.Setenv(envGoghLogLevel, "debug"))
			require.NoError(t, os.Setenv(envGoghLogDate, "0"))
			require.NoError(t, os.Setenv(envGoghLogTime, "0"))
			require.NoError(t, os.Setenv(envGoghLogMicroSeconds, "0"))
			require.NoError(t, os.Setenv(envGoghLogLongFile, "0"))
			require.NoError(t, os.Setenv(envGoghLogShortFile, "0"))
			require.NoError(t, os.Setenv(envGoghLogUTC, "0"))
			require.NoError(t, os.Setenv(envGoghRoot, "/baz:/bux"))

			cfg2, err := GetEnvarConfig()
			require.NoError(t, err)

			cfg := MergeConfig(cfg1, cfg2) // prior after config
			assert.Equal(t, "tokenx2", cfg.GitHubToken())
			assert.Equal(t, "hostx2", cfg.GitHubHost())
			assert.Equal(t, "kyoh87", cfg.GitHubUser())
			assert.Equal(t, "debug", cfg.LogLevel())
			assert.Equal(t, 0, cfg.LogFlags())
			assert.Equal(t, []string{"/baz", "/bux"}, cfg.Root())
			assert.Equal(t, "/baz", cfg.PrimaryRoot())
			assert.Equal(t, os.Stderr, cfg.Stderr())
			assert.Equal(t, os.Stdout, cfg.Stdout())
		})

		t.Run("no overwritten config", func(t *testing.T) {
			resetEnv(t)

			cfg2, err := GetEnvarConfig()
			require.NoError(t, err)

			cfg := MergeConfig(cfg1, cfg2) // prior after config
			assert.Equal(t, "tokenx1", cfg.GitHubToken())
			assert.Equal(t, "hostx1", cfg.GitHubHost())
			assert.Equal(t, "kyoh86", cfg.GitHubUser())
			assert.Equal(t, "trace", cfg.LogLevel())
			assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
			assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
			assert.Equal(t, "/foo", cfg.PrimaryRoot())
			assert.Equal(t, os.Stderr, cfg.Stderr())
			assert.Equal(t, os.Stdout, cfg.Stdout())
		})

	})

	t.Run("get context from envar", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
		require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))
		require.NoError(t, os.Setenv(envGoghLogDate, "1"))
		require.NoError(t, os.Setenv(envGoghLogTime, "1"))
		require.NoError(t, os.Setenv(envGoghLogMicroSeconds, "1"))
		require.NoError(t, os.Setenv(envGoghLogLongFile, "1"))
		require.NoError(t, os.Setenv(envGoghLogShortFile, "1"))
		require.NoError(t, os.Setenv(envGoghLogUTC, "1"))
		require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))
		cfg, err := GetEnvarConfig()
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", cfg.GitHubToken())
		assert.Equal(t, "hostx1", cfg.GitHubHost())
		assert.Equal(t, "kyoh86", cfg.GitHubUser())
		assert.Equal(t, "trace", cfg.LogLevel())
		assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
		assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
		assert.Equal(t, "/foo", cfg.PrimaryRoot())
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
	})

	t.Run("get context from config", func(t *testing.T) {
		resetEnv(t)
		cfg, err := LoadConfig(bytes.NewBufferString(`
root:
- /foo
- /bar

log:
  level: trace
  date: true
  time: true
  microseconds: true
  longfile: true
  shortfile: true
  utc: true

github:
  token: tokenx1
  user: kyoh86
  host: hostx1
`))
		require.NoError(t, err)
		assert.Equal(t, "tokenx1", cfg.GitHubToken())
		assert.Equal(t, "hostx1", cfg.GitHubHost())
		assert.Equal(t, "kyoh86", cfg.GitHubUser())
		assert.Equal(t, "trace", cfg.LogLevel())
		assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
		assert.Equal(t, []string{"/foo", "/bar"}, cfg.Root())
		assert.Equal(t, "/foo", cfg.PrimaryRoot())
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
	})

	t.Run("get default context", func(t *testing.T) {
		resetEnv(t)
		cfg := DefaultConfig()
		assert.Equal(t, "", cfg.GitHubToken())
		assert.Equal(t, "github.com", cfg.GitHubHost())
		assert.Equal(t, "", cfg.GitHubUser())
		assert.Equal(t, "warn", cfg.LogLevel())
		assert.NotEmpty(t, cfg.Root())
		assert.NotEmpty(t, cfg.PrimaryRoot())
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
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

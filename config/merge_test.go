package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeConfig(t *testing.T) {
	resetEnv(t)

	require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx1"))
	require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx1"))
	require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh86"))
	require.NoError(t, os.Setenv(envGoghLogLevel, "trace"))
	require.NoError(t, os.Setenv(envGoghLogDate, "yes"))
	require.NoError(t, os.Setenv(envGoghLogTime, "yes"))
	require.NoError(t, os.Setenv(envGoghLogMicroSeconds, "yes"))
	require.NoError(t, os.Setenv(envGoghLogLongFile, "yes"))
	require.NoError(t, os.Setenv(envGoghLogShortFile, "yes"))
	require.NoError(t, os.Setenv(envGoghLogUTC, "yes"))
	require.NoError(t, os.Setenv(envGoghRoot, "/foo:/bar"))

	cfg1, err := GetEnvarConfig()
	require.NoError(t, err)

	t.Run("full overwritten config", func(t *testing.T) {
		require.NoError(t, os.Setenv(envGoghGitHubToken, "tokenx2"))
		require.NoError(t, os.Setenv(envGoghGitHubHost, "hostx2"))
		require.NoError(t, os.Setenv(envGoghGitHubUser, "kyoh87"))
		require.NoError(t, os.Setenv(envGoghLogLevel, "debug"))
		require.NoError(t, os.Setenv(envGoghLogDate, "no"))
		require.NoError(t, os.Setenv(envGoghLogTime, "no"))
		require.NoError(t, os.Setenv(envGoghLogMicroSeconds, "no"))
		require.NoError(t, os.Setenv(envGoghLogLongFile, "no"))
		require.NoError(t, os.Setenv(envGoghLogShortFile, "no"))
		require.NoError(t, os.Setenv(envGoghLogUTC, "no"))
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

	resetEnv(t)
	assert.Equal(t, EmptyBoolOption, mergeBoolOption(EmptyBoolOption, EmptyBoolOption))
	assert.Equal(t, TrueOption, mergeBoolOption(TrueOption, EmptyBoolOption))
	assert.Equal(t, FalseOption, mergeBoolOption(FalseOption, EmptyBoolOption))
	assert.Equal(t, TrueOption, mergeBoolOption(EmptyBoolOption, TrueOption))
	assert.Equal(t, FalseOption, mergeBoolOption(TrueOption, FalseOption))
	assert.Equal(t, TrueOption, mergeBoolOption(FalseOption, TrueOption))
}

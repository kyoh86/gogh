package env_test

import (
	"os"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// NOTE: these tests include for generators.
	// So it should not be payed attension for coverage.
	t.Run("emptyConfig", func(t *testing.T) {
		config, err := env.LoadConfig(strings.NewReader("{}"))
		require.NoError(t, err)
		assert.Nil(t, config.GithubHost)
		assert.Nil(t, config.Roots)
	})

	t.Run("filledConfig", func(t *testing.T) {
		configRaw := `
githubHost: example.com
roots:
  - foo
  - bar`
		config, err := env.LoadConfig(strings.NewReader(configRaw))
		require.NoError(t, err)
		assert.Equal(t, "example.com", config.GithubHost.Value())
		assert.EqualValues(t, []string{"foo", "bar"}, config.Roots.Value())
	})

	t.Run("emptyCache", func(t *testing.T) {
		cache, err := env.LoadCache(strings.NewReader("{}"))
		require.NoError(t, err)
		assert.Nil(t, cache.GithubUser)
	})

	t.Run("filledCache", func(t *testing.T) {
		cacheRaw := `
githubUser: kyoh86`
		cache, err := env.LoadCache(strings.NewReader(cacheRaw))
		require.NoError(t, err)
		assert.Equal(t, "kyoh86", cache.GithubUser.Value())
	})

	t.Run("emptyEnvar", func(t *testing.T) {
		envar, err := env.LoadEnvar()
		require.NoError(t, err)
		assert.Nil(t, envar.GithubHost)
		assert.Nil(t, envar.Roots)
	})

	t.Run("filledEnvar", func(t *testing.T) {
		os.Setenv("GOGH_GITHUB_TOKEN", "dummy-token")
		os.Setenv("GOGH_GITHUB_HOST", "example.com")
		os.Setenv("GOGH_ROOTS", "foo:bar")
		envar, err := env.LoadEnvar()
		require.NoError(t, err)
		if assert.NotNil(t, envar.GithubToken) {
			assert.Equal(t, "dummy-token", envar.GithubToken.Value())
		}
		if assert.NotNil(t, envar.GithubHost) {
			assert.Equal(t, "example.com", envar.GithubHost.Value())
		}
		if assert.NotNil(t, envar.Roots) {
			assert.EqualValues(t, []string{"foo", "bar"}, envar.Roots.Value())
		}
	})

	t.Run("mergeEmpty", func(t *testing.T) {
		merged := env.Merge(env.Envar{}, env.Cache{}, env.Keyring{}, env.Config{})
		assert.Empty(t, merged.GithubToken())
		assert.Equal(t, "github.com", merged.GithubHost())
		assert.NotEmpty(t, merged.Roots())

	})

	t.Run("mergeOverride", func(t *testing.T) {
		configRaw := `
githubHost: host1
roots:
  - root1a
  - root1b`
		config, err := env.LoadConfig(strings.NewReader(configRaw))
		require.NoError(t, err)

		os.Setenv("GOGH_GITHUB_TOKEN", "dummy-token")
		os.Setenv("GOGH_GITHUB_HOST", "host2")
		os.Setenv("GOGH_ROOTS", "root2a")
		envar, err := env.LoadEnvar()
		require.NoError(t, err)

		merged := env.Merge(envar, env.Cache{}, env.Keyring{}, config)
		assert.Equal(t, "dummy-token", merged.GithubToken())
		assert.Equal(t, "host2", merged.GithubHost())
		assert.EqualValues(t, []string{"root2a"}, merged.Roots())
	})
}

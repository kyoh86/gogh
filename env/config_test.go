package env

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	// NOTE: these tests include for generators.
	t.Run("emptyFile", func(t *testing.T) {
		file, err := loadYAML(strings.NewReader("{}"))
		require.NoError(t, err)
		assert.Nil(t, file.GithubHost)
		assert.Nil(t, file.Roots)
	})

	t.Run("filledFile", func(t *testing.T) {
		fileRaw := `
githubHost: example.com
roots:
  - foo
  - bar`
		file, err := loadYAML(strings.NewReader(fileRaw))
		require.NoError(t, err)
		assert.Equal(t, "example.com", file.GithubHost.Value())
		assert.EqualValues(t, []string{"foo", "bar"}, file.Roots.Value())
	})

	t.Run("emptyEnvar", func(t *testing.T) {
		envar, err := getEnvar("GOGH_")
		require.NoError(t, err)
		assert.Nil(t, envar.GithubHost)
		assert.Nil(t, envar.Roots)
	})

	t.Run("filledEnvar", func(t *testing.T) {
		os.Setenv("GOGH_GITHUB_TOKEN", "dummy-token")
		os.Setenv("GOGH_GITHUB_HOST", "example.com")
		os.Setenv("GOGH_ROOTS", "foo:bar")
		envar, err := getEnvar("GOGH_")
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

	/*
			t.Run("mergeEmpty", func(t *testing.T) {
				merged := mergeAll(YAML{}, Keyring{}, Envar{})
				assert.Empty(t, merged.GithubToken())
				assert.Equal(t, "github.com", merged.GithubHost())
				assert.NotEmpty(t, merged.Roots())

			})

			t.Run("mergeOverride", func(t *testing.T) {
				fileRaw := `
		githubHost: host1
		roots:
		  - root1a
		  - root1b`
				file, err := loadYAML(strings.NewReader(fileRaw))
				require.NoError(t, err)

				os.Setenv("GOGH_GITHUB_TOKEN", "dummy-token")
				os.Setenv("GOGH_GITHUB_HOST", "host2")
				os.Setenv("GOGH_ROOTS", "root2a")
				envar, err := getEnvar("GOGH_")
				require.NoError(t, err)

				merged := mergeAll(file, Keyring{}, envar)
				assert.Equal(t, "dummy-token", merged.GithubToken())
				assert.Equal(t, "host2", merged.GithubHost())
				assert.EqualValues(t, []string{"root2a"}, merged.Roots())
			})
	*/
}

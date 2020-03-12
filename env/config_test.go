package env

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UNDONE: move this test to appenv

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
}

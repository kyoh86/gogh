package repo

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	roots = []string{tmp}

	t.Run("FromFullPath", func(t *testing.T) {
		r, err := FromFullPath(filepath.Join(tmp, "github.com", "kyoh86", "gogh"))
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", r.NonHostPath())
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, r.Subpaths())
	})

	t.Run("FromURL", func(t *testing.T) {
		githubURL, _ := url.Parse("ssh://git@github.com/kyoh86/gogh.git")
		r, err := FromURL(githubURL)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(tmp, "github.com", "kyoh86", "gogh"), r.FullPath)
	})
}

// https://gist.github.com/kyanny/c231f48e5d08b98ff2c3
func TestList_Symlink(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	symDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	roots = []string{root}

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	require.NoError(t, err)

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "gogh"))
	require.NoError(t, err)

	paths := []string{}
	Walk(func(repo *Local) error {
		paths = append(paths, repo.RelPath)
		return nil
	})

	assert.Len(t, paths, 2)
}

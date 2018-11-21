package gogh

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromFullPath(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{roots: []string{tmp}}

	t.Run("in primary root", func(t *testing.T) {
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		r, err := FromFullPath(&ctx, path)
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", r.NonHostPath())
		assert.Equal(t, path, r.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, r.Subpaths())
		assert.True(t, r.IsInPrimaryRoot(&ctx))
		t.Run("TestMatch", func(t *testing.T) {
			assert.True(t, r.Matches("gogh"))
			assert.True(t, r.Matches("kyoh86/gogh"))
			assert.True(t, r.Matches("github.com/kyoh86/gogh"))

			assert.False(t, r.Matches("gigh"))
			assert.False(t, r.Matches("kyoh85/gogh"))
			assert.False(t, r.Matches("githib.com/kyoh86/gogh"))
			assert.False(t, r.Matches("github.com/kyoh86"))
		})
	})
	t.Run("secondary path", func(t *testing.T) {
		tmp2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
		require.NoError(t, err)
		ctx := ctx
		ctx.roots = append(ctx.roots, tmp2)
		path := filepath.Join(tmp2, "github.com", "kyoh86", "gogh")
		r, err := FromFullPath(&ctx, path)
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", r.NonHostPath())
		assert.Equal(t, path, r.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, r.Subpaths())
		assert.False(t, r.IsInPrimaryRoot(&ctx))
	})
	t.Run("not in root path", func(t *testing.T) {
		path := filepath.Join("/src", "github.com", "kyoh86", "gogh")
		_, err := FromFullPath(&ctx, path)
		assert.Error(t, err, "no repository found for: "+path)
	})

}

func TestFromURL(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{roots: []string{tmp}}

	path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
	t.Run("not existing repository", func(t *testing.T) {
		r, err := FromURL(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, r.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, r.Subpaths())
	})
	t.Run("not supported host URL", func(t *testing.T) {
		_, err := FromURL(&ctx, parseURL(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.Error(t, err, `not supported host: "example.com"`)
	})
	t.Run("existing repository", func(t *testing.T) {
		// Create dummy repository
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, "github.com", "kyoh85", "gogh", ".git"), 0755))
		// Create target repository
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))
		defer func() {
			require.NoError(t, os.RemoveAll(path))
		}()
		r, err := FromURL(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, r.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, r.Subpaths())
	})
}

func parseURL(t *testing.T, text string) *url.URL {
	t.Helper()
	u, err := url.Parse(text)
	require.NoError(t, err)
	return u
}

func TestCheckURL(t *testing.T) {
	t.Run("valid github URL", func(t *testing.T) {
		ctx := implContext{}
		assert.NoError(t, CheckURL(&ctx, parseURL(t, "https://github.com/kyoh86/gogh")))
	})
	t.Run("valid GHE URL", func(t *testing.T) {
		ctx := implContext{
			gheHosts: []string{"example.com"},
		}
		assert.NoError(t, CheckURL(&ctx, parseURL(t, "https://example.com/kyoh86/gogh")))
	})
	t.Run("valid github URL with trailing slashes", func(t *testing.T) {
		ctx := implContext{}
		assert.NoError(t, CheckURL(&ctx, parseURL(t, "https://github.com/kyoh86/gogh/")))
	})
	t.Run("valid GHE URL with trailing slashes", func(t *testing.T) {
		ctx := implContext{
			gheHosts: []string{"example.com"},
		}
		assert.NoError(t, CheckURL(&ctx, parseURL(t, "https://example.com/kyoh86/gogh/")))
	})
	t.Run("not supported host URL", func(t *testing.T) {
		ctx := implContext{
			gheHosts: []string{"example.com"},
		}
		assert.Error(t, CheckURL(&ctx, parseURL(t, "https://kyoh86.work/kyoh86/gogh")), `not supported host: "kyoh86.work"`)
	})
	t.Run("invalid path on GitHub", func(t *testing.T) {
		ctx := implContext{}
		assert.Error(t, CheckURL(&ctx, parseURL(t, "https://github.com/kyoh86/gogh/blob/master/README.md")), `URL should be formed 'schema://hostname/user/name'`)
	})
	t.Run("invalid path on GHE", func(t *testing.T) {
		ctx := implContext{
			gheHosts: []string{"example.com"},
		}
		assert.Error(t, CheckURL(&ctx, parseURL(t, "https://example.com/kyoh86/gogh/blob/master/README.md")), `URL should be formed 'schema://hostname/user/name'`)
	})
}

func TestWalk(t *testing.T) {
	neverCalled := func(t *testing.T) func(*LocalRepo) error {
		return func(*LocalRepo) error {
			t.Fatal("should not be called but...")
			return nil
		}
	}
	t.Run("Not existing root", func(t *testing.T) {
		t.Run("primary root", func(t *testing.T) {
			ctx := implContext{roots: []string{"/that/will/never/exist"}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
		t.Run("secondary root", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			ctx := implContext{roots: []string{tmp, "/that/will/never/exist"}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
	})

	t.Run("Root specifies a file", func(t *testing.T) {
		t.Run("Primary root is a file", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
			ctx := implContext{roots: []string{filepath.Join(tmp, "foo")}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
		t.Run("Secondary root is a file", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
			ctx := implContext{roots: []string{tmp, filepath.Join(tmp, "foo")}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
	})

	t.Run("through error from callback", func(t *testing.T) {
		tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
		require.NoError(t, err)
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))

		require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
		ctx := implContext{roots: []string{tmp, filepath.Join(tmp, "foo")}}
		err = errors.New("sample error")
		assert.Error(t, err, Walk(&ctx, func(l *LocalRepo) error {
			assert.Equal(t, path, l.FullPath)
			return err
		}))
	})
}

// https://gist.github.com/kyanny/c231f48e5d08b98ff2c3
func TestList_Symlink(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	symDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	ctx := &implContext{roots: []string{root}}

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	require.NoError(t, err)

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "gogh"))
	require.NoError(t, err)

	paths := []string{}
	Walk(ctx, func(repo *LocalRepo) error {
		paths = append(paths, repo.RelPath)
		return nil
	})

	assert.Len(t, paths, 2)
}

package gogh

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLocal(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{roots: []string{tmp}}

	t.Run("in primary root", func(t *testing.T) {
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		l, err := parseLocal(tmp, path)
		require.NoError(t, err)
		assert.Equal(t, path, l.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, l.Subpaths())
		assert.True(t, l.IsInPrimaryRoot(&ctx))
	})
	t.Run("secondary path", func(t *testing.T) {
		tmp2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
		require.NoError(t, err)
		ctx := ctx
		ctx.roots = append(ctx.roots, tmp2)
		path := filepath.Join(tmp2, "github.com", "kyoh86", "gogh")
		l, err := parseLocal(tmp2, path)
		require.NoError(t, err)
		assert.Equal(t, path, l.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, l.Subpaths())
		assert.False(t, l.IsInPrimaryRoot(&ctx))
	})
}

func TestFindLocal(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{roots: []string{tmp}, userName: "kyoh86"}

	path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
	t.Run("not existing repository", func(t *testing.T) {
		l, err := FindLocal(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, l.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, l.Subpaths())
	})
	t.Run("not supported host URL", func(t *testing.T) {
		_, err := FindLocal(&ctx, parseURL(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.Error(t, err, `not supported host: "example.com"`)
	})
	t.Run("existing repository", func(t *testing.T) {
		// Create same name repository
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, "github.com", "kyoh85", "gogh", ".git"), 0755))
		// Create different name repository
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, "github.com", "kyoh86", "foo", ".git"), 0755))
		// Create target repository
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))
		defer func() {
			require.NoError(t, os.RemoveAll(path))
		}()

		t.Run("full name", func(t *testing.T) {
			gotPath, err := FindLocalPath(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
			require.NoError(t, err)
			assert.Equal(t, path, gotPath)
		})

		t.Run("shortest precise name (owner and name)", func(t *testing.T) {
			l, err := FindLocal(&ctx, parseURL(t, "kyoh86/gogh"))
			require.NoError(t, err)
			assert.Equal(t, path, l.FullPath)
			assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, l.Subpaths())
		})

		t.Run("shortest pricese name (name only)", func(t *testing.T) {
			l, err := FindLocal(&ctx, parseURL(t, "foo"))
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(tmp, "github.com", "kyoh86", "foo"), l.FullPath)
			assert.Equal(t, []string{"foo", "kyoh86/foo", "github.com/kyoh86/foo"}, l.Subpaths())
		})
	})
}

func parseURL(t *testing.T, text string) *Remote {
	t.Helper()
	u, err := ParseRemote(text)
	require.NoError(t, err)
	return u
}

func TestWalk(t *testing.T) {
	neverCalled := func(t *testing.T) func(*Local) error {
		return func(*Local) error {
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
			require.NoError(t, WalkInPrimary(&ctx, neverCalled(t)))
		})
		t.Run("Secondary root is a file", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
			ctx := implContext{roots: []string{tmp, filepath.Join(tmp, "foo")}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
			require.NoError(t, WalkInPrimary(&ctx, neverCalled(t)))
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
		assert.Error(t, err, Walk(&ctx, func(l *Local) error {
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
	Walk(ctx, func(l *Local) error {
		paths = append(paths, l.RelPath)
		return nil
	})

	assert.Len(t, paths, 2)
}

func TestQuery(t *testing.T) {
	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)
	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)
	path1 := filepath.Join(root1, "github.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path1, ".git"), 0755))
	path2 := filepath.Join(root1, "github.com", "kyoh85", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path2, ".git"), 0755))
	path3 := filepath.Join(root1, "github.com", "kyoh86", "foo")
	require.NoError(t, os.MkdirAll(filepath.Join(path3, ".git"), 0755))
	path4 := filepath.Join(root1, "example.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path4, ".git"), 0755))
	path5 := filepath.Join(root2, "github.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path5, ".git"), 0755))

	ctx := implContext{roots: []string{root1, root2}}

	assert.NoError(t, Query(&ctx, "never found", Walk, func(*Local) error {
		t.Fatal("should not be called but...")
		return nil
	}))

	t.Run("NameOnly", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
			path4: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "gogh", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
			path4: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "gog", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("OwnerAndName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path4: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "kyoh86/gogh", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialOwnerAndName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path4: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "yoh86/gog", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("FullRemoteName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "github.com/kyoh86/gogh", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialFullRemoteName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "ithub.com/kyoh86/gog", Walk, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("WalkInPrimary", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
			path4: {},
		}
		assert.NoError(t, Query(&ctx, "gogh", WalkInPrimary, func(l *Local) error {
			assert.Contains(t, expect, l.FullPath)
			delete(expect, l.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
}

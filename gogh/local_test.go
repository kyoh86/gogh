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

func TestParseProject(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{root: []string{tmp}, gitHubHost: "github.com"}

	t.Run("in primary root", func(t *testing.T) {
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		p, err := parseProject(&ctx, tmp, path)
		require.NoError(t, err)
		assert.Equal(t, path, p.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
		assert.True(t, p.IsInPrimaryRoot(&ctx))
	})

	t.Run("secondary path", func(t *testing.T) {
		tmp2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
		require.NoError(t, err)
		ctx := ctx
		ctx.root = append(ctx.root, tmp2)
		path := filepath.Join(tmp2, "github.com", "kyoh86", "gogh")
		p, err := parseProject(&ctx, tmp2, path)
		require.NoError(t, err)
		assert.Equal(t, path, p.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
		assert.False(t, p.IsInPrimaryRoot(&ctx))
	})

	t.Run("expect to fail to parse unsupported depth", func(t *testing.T) {
		r, err := parseProject(&ctx, tmp, filepath.Join(tmp, "github.com/kyoh86/gogh/gogh"))
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("expect to fail to parse unsupported host", func(t *testing.T) {
		r, err := parseProject(&ctx, tmp, filepath.Join(tmp, "example.com/kyoh86/gogh"))
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("expect to fail to parse owner name that starts with hyphen", func(t *testing.T) {
		r, err := parseProject(&ctx, tmp, filepath.Join(tmp, "github.com/-kyoh86/gogh"))
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

	t.Run("expect to fail to parse project name that contains invalid charactor", func(t *testing.T) {
		r, err := parseProject(&ctx, tmp, filepath.Join(tmp, "github.com/kyoh86/foo,bar"))
		assert.NotNil(t, err)
		assert.Nil(t, r)
	})

}

func TestFindOrNewProject(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := implContext{root: []string{tmp}, gitHubUser: "kyoh86", gitHubHost: "github.com"}

	path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")

	t.Run("not existing repository", func(t *testing.T) {
		p, err := FindOrNewProject(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, p.FullPath)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
	})
	t.Run("not existing repository with FindProject", func(t *testing.T) {
		_, err := FindProject(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, "project not found")
	})
	t.Run("not supported host URL by FindProject", func(t *testing.T) {
		_, err := FindOrNewProject(&ctx, parseURL(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, `not supported host: "example.com"`)
	})
	t.Run("not supported host URL by FindProjectPath", func(t *testing.T) {
		_, err := FindProjectPath(&ctx, parseURL(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, `not supported host: "example.com"`)
	})
	t.Run("not supported host URL by NewProject", func(t *testing.T) {
		_, err := NewProject(&ctx, parseURL(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, `not supported host: "example.com"`)
	})
	t.Run("fail with invalid root", func(t *testing.T) {
		ctx := implContext{root: []string{"/\x00"}, gitHubUser: "kyoh86", gitHubHost: "github.com"}
		_, err := FindOrNewProject(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
		assert.Error(t, err)
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
			gotPath, err := FindProjectPath(&ctx, parseURL(t, "ssh://git@github.com/kyoh86/gogh.git"))
			require.NoError(t, err)
			assert.Equal(t, path, gotPath)
		})

		t.Run("shortest precise name (owner and name)", func(t *testing.T) {
			p, err := FindOrNewProject(&ctx, parseURL(t, "kyoh86/gogh"))
			require.NoError(t, err)
			assert.Equal(t, path, p.FullPath)
			assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
		})

		t.Run("shortest pricese name (name only)", func(t *testing.T) {
			p, err := FindOrNewProject(&ctx, parseURL(t, "foo"))
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(tmp, "github.com", "kyoh86", "foo"), p.FullPath)
			assert.Equal(t, []string{"foo", "kyoh86/foo", "github.com/kyoh86/foo"}, p.Subpaths())
		})
	})
}

func parseURL(t *testing.T, text string) *Repo {
	t.Helper()
	u, err := ParseRepo(text)
	require.NoError(t, err)
	return u
}

func TestWalk(t *testing.T) {
	neverCalled := func(t *testing.T) func(*Project) error {
		return func(*Project) error {
			t.Fatal("should not be called but...")
			return nil
		}
	}
	t.Run("Not existing root", func(t *testing.T) {
		t.Run("primary root", func(t *testing.T) {
			ctx := implContext{root: []string{"/that/will/never/exist"}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
		t.Run("secondary root", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			ctx := implContext{root: []string{tmp, "/that/will/never/exist"}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
		})
	})

	t.Run("Root specifies a file", func(t *testing.T) {
		t.Run("Primary root is a file", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
			ctx := implContext{root: []string{filepath.Join(tmp, "foo")}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
			require.NoError(t, WalkInPrimary(&ctx, neverCalled(t)))
		})
		t.Run("Secondary root is a file", func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
			ctx := implContext{root: []string{tmp, filepath.Join(tmp, "foo")}}
			require.NoError(t, Walk(&ctx, neverCalled(t)))
			require.NoError(t, WalkInPrimary(&ctx, neverCalled(t)))
		})
	})

	t.Run("through error with invalid project name", func(t *testing.T) {
		tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
		require.NoError(t, err)
		path := filepath.Join(tmp, "github.com", "kyoh--86", "gogh")
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))

		ctx := implContext{root: []string{tmp, filepath.Join(tmp)}}
		assert.NoError(t, Walk(&ctx, neverCalled(t)))
	})

	t.Run("through error from callback", func(t *testing.T) {
		tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
		require.NoError(t, err)
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))

		require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))
		ctx := implContext{root: []string{tmp, filepath.Join(tmp, "foo")}, gitHubHost: "github.com"}
		err = errors.New("sample error")
		assert.EqualError(t, Walk(&ctx, func(p *Project) error {
			assert.Equal(t, path, p.FullPath)
			return err
		}), "sample error")
	})
}

// https://gist.github.com/kyanny/c231f48e5d08b98ff2c3
func TestList_Symlink(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	symDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	ctx := &implContext{root: []string{root}, gitHubHost: "github.com"}

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	require.NoError(t, err)

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "gogh"))
	require.NoError(t, err)

	paths := []string{}
	require.NoError(t, Walk(ctx, func(p *Project) error {
		paths = append(paths, p.RelPath)
		return nil
	}))

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
	path5 := filepath.Join(root2, "github.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path5, ".git"), 0755))

	ctx := implContext{root: []string{root1, root2}, gitHubHost: "github.com"}

	assert.NoError(t, Query(&ctx, "never found", Walk, func(*Project) error {
		t.Fatal("should not be called but...")
		return nil
	}))

	t.Run("NameOnly", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "gogh", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "gog", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("OwnerAndName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "kyoh86/gogh", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialOwnerAndName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "yoh86/gog", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("FullRepoName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "github.com/kyoh86/gogh", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("PartialFullRepoName", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path5: {},
		}
		assert.NoError(t, Query(&ctx, "ithub.com/kyoh86/gog", Walk, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
	t.Run("WalkInPrimary", func(t *testing.T) {
		expect := map[string]struct{}{
			path1: {},
			path2: {},
		}
		assert.NoError(t, Query(&ctx, "gogh", WalkInPrimary, func(p *Project) error {
			assert.Contains(t, expect, p.FullPath)
			delete(expect, p.FullPath)
			return nil
		}))
		assert.Empty(t, expect)
	})
}

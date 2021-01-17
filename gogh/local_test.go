package gogh_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)

	t.Run("success", func(t *testing.T) {
		var spec gogh.RepoSpec
		require.NoError(t, spec.Set("https://github.com/kyoh86/gogh"))
		proj, err := gogh.NewProject(svc.ev, &spec)
		assert.NotNil(t, proj)
		assert.NoError(t, err)
	})

	t.Run("not supported host", func(t *testing.T) {
		var spec gogh.RepoSpec
		require.NoError(t, spec.Set("https://example.com/kyoh86/gogh"))
		proj, err := gogh.NewProject(svc.ev, &spec)
		assert.Error(t, err)
		assert.Nil(t, proj)
	})
}

func TestFindOrNewProject(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)

	path := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")

	t.Run("not existing repository", func(t *testing.T) {
		p, _, err := gogh.FindOrNewProject(svc.ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, p.FullPath)
		assert.False(t, p.Exists)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
	})
	t.Run("not existing repository (in primary)", func(t *testing.T) {
		// Create same name repository in other root
		inOther := filepath.Join(svc.root2, "github.com", "kyoh86", "gogh", ".git")
		require.NoError(t, os.MkdirAll(inOther, 0755))
		defer os.RemoveAll(inOther)
		p, _, err := gogh.FindOrNewProjectInPrimary(svc.ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
		require.NoError(t, err)
		assert.Equal(t, path, p.FullPath)
		assert.False(t, p.Exists)
		assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
	})
	t.Run("not existing repository with FindProject", func(t *testing.T) {
		_, _, err := gogh.FindProject(svc.ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, "project not found")
	})
	t.Run("not existing repository with FindProjectInPrimary", func(t *testing.T) {
		inOther := filepath.Join(svc.root2, "github.com", "kyoh86", "gogh", ".git")
		require.NoError(t, os.MkdirAll(inOther, 0755))
		defer os.RemoveAll(inOther)
		_, _, err := gogh.FindProjectInPrimary(svc.ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
		assert.EqualError(t, err, "project not found")
	})
	t.Run("fail with invalid root", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ev := NewMockEnv(ctrl)
		ev.EXPECT().GithubHost().AnyTimes().Return("github.com")
		ev.EXPECT().Roots().AnyTimes().Return([]string{"/\x00"})

		_, _, err := gogh.FindOrNewProject(ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
		assert.Error(t, err)
	})
	t.Run("not supported host", func(t *testing.T) {
		_, _, err := gogh.FindOrNewProject(svc.ev, mustParseRepoSpec(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.Error(t, err)
	})
	t.Run("finding path with not supported host", func(t *testing.T) {
		_, err := gogh.FindProjectPath(svc.ev, mustParseRepoSpec(t, "ssh://git@example.com/kyoh86/gogh.git"))
		assert.Error(t, err)
	})
	t.Run("existing repository", func(t *testing.T) {
		// Create same name repository
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root1, "github.com", "kyoh85", "gogh", ".git"), 0755))
		// Create different name repository
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root1, "github.com", "kyoh86", "foo", ".git"), 0755))
		// Create target repository
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))
		defer func() {
			require.NoError(t, os.RemoveAll(path))
		}()

		t.Run("full name", func(t *testing.T) {
			gotPath, err := gogh.FindProjectPath(svc.ev, mustParseRepoSpec(t, "ssh://git@github.com/kyoh86/gogh.git"))
			require.NoError(t, err)
			assert.Equal(t, path, gotPath)
		})

		t.Run("shortest precise name (owner and name)", func(t *testing.T) {
			p, _, err := gogh.FindOrNewProject(svc.ev, mustParseRepoSpec(t, "kyoh86/gogh"))
			require.NoError(t, err)
			assert.Equal(t, path, p.FullPath)
			assert.Equal(t, []string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}, p.Subpaths())
		})

		t.Run("shortest pricese name (name only)", func(t *testing.T) {
			p, _, err := gogh.FindOrNewProject(svc.ev, mustParseRepoSpec(t, "foo"))
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(svc.root1, "github.com", "kyoh86", "foo"), p.FullPath)
			assert.Equal(t, []string{"foo", "kyoh86/foo", "github.com/kyoh86/foo"}, p.Subpaths())
		})
	})
}

func TestWalk(t *testing.T) {
	neverCalled := func(t *testing.T) func(*gogh.Project) error {
		return func(*gogh.Project) error {
			t.Fatal("should not be called but...")
			return nil
		}
	}

	t.Run("empty roots", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		assert.NoError(t, gogh.Walk(svc.ev, neverCalled(t)))
	})

	t.Run("unexist roots", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)

		configCtrl := gomock.NewController(t)
		ev := NewMockEnv(configCtrl)
		ev.EXPECT().GithubUser().AnyTimes().Return("kyoh86")
		ev.EXPECT().GithubHost().AnyTimes().Return("github.com")
		ev.EXPECT().Roots().AnyTimes().Return([]string{filepath.Join(svc.root1, "unexist"), svc.root2})
		defer configCtrl.Finish()

		assert.NoError(t, gogh.Walk(ev, neverCalled(t)))
	})

	t.Run("cover all projects", func(t *testing.T) {
		svc := initTest(t)
		defer svc.teardown(t)
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root1, "github.com", "kyoh86", "proj-1", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root1, "github.com", "kyoh86", "proj-2", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "github.com", "kyoh86", "proj-3", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "github.com", "kyoh86", "proj-4", ".git"), os.ModePerm))

		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "invalid.com", "kyoh86", "invalid-1", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "github.com", "kyoh86", "too-deep", "invalid-2", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "github.com", "in'valid", "invalid-3", ".git"), os.ModePerm))
		require.NoError(t, os.MkdirAll(filepath.Join(svc.root2, "github.com", "kyoh86", "in'valid-4", ".git"), os.ModePerm))

		var (
			found1 bool
			found2 bool
			found3 bool
			found4 bool
		)
		require.NoError(t, gogh.Walk(svc.ev, func(p *gogh.Project) error {
			switch p.RelPath {
			case filepath.ToSlash(filepath.Join("github.com", "kyoh86", "proj-1")):
				assert.True(t, p.IsInPrimaryRoot(svc.ev))
				found1 = true
			case filepath.ToSlash(filepath.Join("github.com", "kyoh86", "proj-2")):
				assert.True(t, p.IsInPrimaryRoot(svc.ev))
				found2 = true
			case filepath.ToSlash(filepath.Join("github.com", "kyoh86", "proj-3")):
				assert.False(t, p.IsInPrimaryRoot(svc.ev))
				found3 = true
			case filepath.ToSlash(filepath.Join("github.com", "kyoh86", "proj-4")):
				assert.False(t, p.IsInPrimaryRoot(svc.ev))
				found4 = true
			default:
				return fmt.Errorf("not prepared project %s", p.RelPath)
			}
			return nil
		}))
		assert.True(t, found1)
		assert.True(t, found2)
		assert.True(t, found3)
		assert.True(t, found4)
	})

	t.Run("through error with invalid project name", func(t *testing.T) {
		tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
		require.NoError(t, err)
		path := filepath.Join(tmp, "github.com", "kyoh--86", "gogh")
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ev := NewMockEnv(ctrl)
		ev.EXPECT().Roots().AnyTimes().Return([]string{tmp, filepath.Join(tmp)})
		ev.EXPECT().GithubHost().AnyTimes().Return("github.com")

		assert.NoError(t, gogh.Walk(ev, neverCalled(t)))
	})

	t.Run("through error from callback", func(t *testing.T) {
		tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
		require.NoError(t, err)
		path := filepath.Join(tmp, "github.com", "kyoh86", "gogh")
		require.NoError(t, os.MkdirAll(filepath.Join(path, ".git"), 0755))
		require.NoError(t, ioutil.WriteFile(filepath.Join(tmp, "foo"), nil, 0644))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ev := NewMockEnv(ctrl)
		ev.EXPECT().GithubHost().AnyTimes().Return("github.com")
		ev.EXPECT().Roots().AnyTimes().Return([]string{tmp, filepath.Join(tmp, "foo")})

		err = errors.New("sample error")
		assert.EqualError(t, gogh.Walk(ev, func(p *gogh.Project) error {
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ev := NewMockEnv(ctrl)
	ev.EXPECT().GithubHost().AnyTimes().Return("github.com")
	ev.EXPECT().Roots().AnyTimes().Return([]string{root})

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	require.NoError(t, err)

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "gogh"))
	require.NoError(t, err)

	var lock sync.Mutex
	paths := []string{}
	require.NoError(t, gogh.Walk(ev, func(p *gogh.Project) error {
		lock.Lock()
		defer lock.Unlock()
		paths = append(paths, p.RelPath)
		return nil
	}))

	assert.Len(t, paths, 2)
}

func TestQuery(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)

	path1 := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path1, ".git"), 0755))
	path2 := filepath.Join(svc.root1, "github.com", "kyoh85", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path2, ".git"), 0755))
	path3 := filepath.Join(svc.root1, "github.com", "kyoh86", "foo")
	require.NoError(t, os.MkdirAll(filepath.Join(path3, ".git"), 0755))
	path5 := filepath.Join(svc.root2, "github.com", "kyoh86", "gogh")
	require.NoError(t, os.MkdirAll(filepath.Join(path5, ".git"), 0755))

	assert.NoError(t, gogh.Query(svc.ev, "never found", gogh.Walk, func(*gogh.Project) error {
		t.Fatal("should not be called but...")
		return nil
	}))

	t.Run("NameOnly", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path2, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "gogh", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("PartialName", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path2, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "gog", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("OwnerAndName", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "kyoh86/gogh", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("PartialOwnerAndName", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "yoh86/gog", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("FullRepoName", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "github.com/kyoh86/gogh", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("PartialFullRepoName", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path5, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "ithub.com/kyoh86/gog", gogh.Walk, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
	t.Run("WalkInPrimary", func(t *testing.T) {
		var expect sync.Map
		expect.Store(path1, struct{}{})
		expect.Store(path2, struct{}{})
		assert.NoError(t, gogh.Query(svc.ev, "gogh", gogh.WalkInPrimary, func(p *gogh.Project) error {
			_, ok := expect.Load(p.FullPath)
			assert.True(t, ok, p.FullPath)
			expect.Delete(p.FullPath)
			return nil
		}))
		expect.Range(func(key interface{}, value interface{}) bool {
			assert.Failf(t, "not found by walking", "%#v is not found by walking", key)
			return true
		})
	})
}
